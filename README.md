# AI Ассистенты

Fullstack-приложение — каталог AI-ассистентов с ролевой моделью, историей запусков и мок-интеграцией LLM.

## Запуск

```bash
docker-compose up --build
```

- Backend: http://localhost:8080
- Frontend: http://localhost:3000
- Healthcheck: `GET http://localhost:8080/_info` → `200 OK`

Все переменные окружения имеют дефолтные значения в `docker-compose.yaml` — проект запускается без ручной настройки.

## Запуск тестов

### Backend — юнит-тесты

```bash
cd backend
go test ./internal/usecase/... -v
```

### Backend — coverage

```bash
go test ./internal/usecase/... -coverprofile=coverage.out
go tool cover -func=coverage.out
```

### Backend — E2E-тест

Проверяет сквозной сценарий: admin создаёт категорию → admin создаёт ассистента → user запускает ассистента. Также проверяет фильтрацию ассистентов по `categoryId` и запрет запуска неактивного ассистента.

Требует запущенного Postgres:

```bash
TEST_DSN="postgres://postgres:postgres@localhost:5432/ai_assistants?sslmode=disable" \
  go test ./internal/usecase/... -tags=e2e -v -run TestE2E
```

### Frontend

```bash
cd frontend
npm test              # watch-режим
npm run test:coverage # с отчётом coverage
```

## Покрытие тестами (Coverage)

### Backend

| Функция | Покрытие |
|---|---|
| `CategoryUsecase.Create` | 100% |
| `AssistantUsecase.Create` | 87.5% |
| `AssistantUsecase.List` | 66.7% |
| `RunUsecase.Run` | 91.3% |
| `auth.*` | 0% (не тестировался юнит-тестами) |
| **Итого по usecase** | **31.9%** |

Итоговый процент низкий, потому что в подсчёт попадают все файлы пакета, включая `AuthUsecase` с нулевым покрытием. По ключевым бизнес-сценариям покрытие высокое: создание категории, создание ассистента (включая негативные кейсы — пустой системный промпт, несуществующая категория), запуск ассистента (успех, неактивный ассистент, ошибка LLM), фильтрация по `categoryId`. `AuthUsecase` не покрыт юнит-тестами — логика аутентификации проверяется через E2E-сценарий.

### Frontend

| Файл | Statements | Branches | Functions |
|---|---|---|---|
| `StatusBadge` | 100% | 100% | 100% |
| `Button` | 100% | 77.8% | 100% |
| `Checkbox` | 100% | 100% | 100% |
| `Input` | 100% | 71.4% | 100% |
| `NewAssistant` | 97.5% | 86.4% | 100% |
| `Login` | 45.8% | 68% | 72.7% |
| **Итого по всем файлам** | **18.83%** | **55.14%** | **36.53%** |

Общий показатель низкий из-за непокрытых страниц (`AssistantDetail`, `Assistants`, `MyRuns`, `AllRuns`, `EditAssistant`) и инфраструктурных файлов (роутер, стор, API-слой). Все ключевые UI-компоненты и критичные формы покрыты хорошо.

## Схема БД и индексы

```
users
  id UUID PK, email VARCHAR UNIQUE, password VARCHAR, role VARCHAR, created_at

categories
  id UUID PK, name VARCHAR UNIQUE, description TEXT, created_at

assistants
  id UUID PK, category_id UUID FK(categories), name VARCHAR,
  description TEXT, model VARCHAR, system_prompt TEXT,
  example_user_prompt TEXT, is_active BOOLEAN, created_at, updated_at

runs
  id UUID PK, assistant_id UUID FK(assistants), user_id UUID FK(users),
  model VARCHAR, user_prompt TEXT, output TEXT,
  status VARCHAR(pending|success|failed), error TEXT, created_at
```

Индексы:
- `assistants(category_id)` — фильтрация по категории
- `assistants(name)`, `assistants(description)` — поиск по названию и описанию
- `assistants(is_active)` — фильтрация активных
- `runs(user_id)` — история запусков пользователя
- `runs(assistant_id)` — фильтрация запусков по ассистенту
- `runs(status)` — фильтрация по статусу
- `runs(created_at DESC)` — сортировка по дате

## Обработка запуска ассистента

Порядок действий при `POST /assistants/{id}/run`:

1. Загрузить ассистента по ID — если не найден, вернуть ошибку
2. Проверить `isActive` — если неактивен, вернуть ошибку без создания записи
3. Создать запись `run` со статусом `pending`
4. Вызвать LLM-провайдер с системным промптом ассистента и пользовательским промптом
5. При успехе — обновить `run`: статус `success`, записать `output`
6. При ошибке — обновить `run`: статус `failed`, записать `error`
7. Вернуть итоговый `run` клиенту

Запись сохраняется в БД в любом случае — и при успехе, и при ошибке LLM.

## LLM-провайдер

### Мок-провайдер

По умолчанию (`LLM_PROVIDER=mock`) используется детерминированный мок из `internal/llm/mock.go`. Принимает на вход:

```go
type Request struct {
    Model        string
    SystemPrompt string
    UserPrompt   string
}
```

Возвращает строку, которая детерминированно включает модель, системный промпт и пользовательский запрос. Ответ всегда одинаков для одинаковых входных данных — тесты стабильны и не зависят от внешних сервисов.

Пример итогового запроса, который был бы отправлен в реальный LLM:

```json
{
  "model": "gpt-4o-mini",
  "messages": [
    {
      "role": "system",
      "content": "Составь простой домашний рецепт из ингредиентов пользователя: дай название, примерное время и 3-5 коротких шагов."
    },
    {
      "role": "user",
      "content": "курица, лаваш, огурцы, томаты, соус"
    }
  ]
}
```

### Обработка ошибок провайдера

- **Ошибка провайдера** — `run` сохраняется со статусом `failed`, поле `error` содержит описание ошибки, API возвращает ошибку клиенту
- **Таймаут** — обрабатывается как ошибка провайдера, `run` сохраняется со статусом `failed`
- **Невалидный ответ** — если провайдер вернул пустой output, `run` сохраняется со статусом `failed`


## Стейт-менеджмент на фронтенде

Используется **Zustand** — минимальная библиотека без бойлерплейта. Хранится три стора:

- `auth.js` — токен и данные пользователя (сохраняются в `localStorage`)
- `theme.js` — тема (сохраняется в `localStorage`)
- `runs.js` — `lastRunAt`: timestamp последнего запуска, используется как сигнал для обновления истории запусков в `MyRuns`

Redux избыточен для задачи такого масштаба. TanStack Query не использовался для запусков, чтобы не усложнять стек — обновление через `lastRunAt` в Zustand решает задачу инвалидации просто и явно.

## Системный промпт

Системный промпт является внутренней настройкой ассистента и не передаётся обычным пользователям. Фронтенд отправляет при запуске только пользовательский промпт — системный промпт подставляет бэкенд самостоятельно из БД.

На бэкенде в `toAssistantResponse` поле `systemPrompt` включается в ответ только если запрос выполняет администратор. Для обычного пользователя поле отсутствует в JSON (`omitempty`). Фронтенд его не получает и не отображает.

## Принятые самостоятельные решения

- **Фильтры каталога синхронизированы с query-параметрами** — при обновлении страницы или копировании ссылки состояние фильтров сохраняется
- **JWT хранится в `localStorage`** через Zustand-стор — единственное место хранения токена в приложении
- **Миграции применяются автоматически** при старте backend-сервиса через `golang-migrate`
- **`/dummyLogin` возвращает фиксированные UUID** для каждой роли, заданные через переменные окружения `ADMIN_UUID` и `USER_UUID` — это обеспечивает воспроизводимость сценариев с историей запусков
- **Пагинация реализована на бэкенде** — фронтенд передаёт `page` и `pageSize`, получает `total` для отрисовки пагинатора

## Дополнительно реализовано

- Регистрация и авторизация по email/паролю (bcrypt)
- Тёмная тема с сохранением в `localStorage`
