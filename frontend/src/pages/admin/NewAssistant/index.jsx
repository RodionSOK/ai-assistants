import { useState, useEffect } from "react";
import { useNavigate, Link } from "react-router-dom";
import { createAssistant } from "@/api/assistants";
import { getCategories } from "@/api/categories";
import Input from "@/components/ui/Input";
import Textarea from "@/components/ui/Textarea";
import Select from "@/components/ui/Select";
import Button from "@/components/ui/Button";
import "./NewAssistant.css";

export default function NewAssistant() {
  const navigate = useNavigate();

  const [categories, setCategories] = useState([]);
  const [form, setForm] = useState({
    categoryId: "",
    name: "",
    description: "",
    model: "",
    systemPrompt: "",
    exampleUserPrompt: "",
    isActive: true,
  });
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);
  const [serverError, setServerError] = useState(null);

  useEffect(() => {
    getCategories()
      .then(({ categories }) => setCategories(categories))
      .catch(() => {});
  }, []);

  const set = (key) => (e) => {
    const value = e.target.type === "checkbox" ? e.target.checked : e.target.value;
    setForm((prev) => ({ ...prev, [key]: value }));
    setErrors((prev) => ({ ...prev, [key]: undefined }));
  };

  const validate = () => {
    const e = {};
    if (!form.categoryId) e.categoryId = "Выберите категорию";
    if (!form.name.trim()) e.name = "Название обязательно";
    if (!form.description.trim()) e.description = "Описание обязательно";
    if (!form.model.trim()) e.model = "Модель обязательна";
    if (!form.systemPrompt.trim()) e.systemPrompt = "Системный промпт обязателен";
    return e;
  };

  const handleSubmit = async () => {
    const e = validate();
    if (Object.keys(e).length) {
      setErrors(e);
      return;
    }

    setLoading(true);
    setServerError(null);

    try {
      await createAssistant({
        ...form,
        name: form.name.trim(),
        description: form.description.trim(),
        model: form.model.trim(),
        systemPrompt: form.systemPrompt.trim(),
        exampleUserPrompt: form.exampleUserPrompt.trim() || null,
      });
      navigate("/assistants");
    } catch (e) {
      const msg = e.response?.data?.error?.message;
      setServerError(msg || "Не удалось создать ассистента");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="new-assistant">
      <div className="new-assistant__back">
        <Link to="/assistants">← Назад</Link>
      </div>

      <div className="new-assistant__card">
        <h1 className="new-assistant__title">Новый ассистент</h1>

        <Select
          label="Категория"
          placeholder="Выберите категорию..."
          value={form.categoryId}
          onChange={set("categoryId")}
          error={errors.categoryId}
          options={categories.map((c) => ({ value: c.id, label: c.name }))}
        />

        <Input
          label="Название"
          placeholder="Например: Повар"
          value={form.name}
          onChange={set("name")}
          error={errors.name}
        />

        <Textarea
          label="Описание"
          placeholder="Чем занимается ассистент"
          value={form.description}
          rows={2}
          onChange={set("description")}
          error={errors.description}
        />

        <Input
          label="Модель"
          placeholder="Например: mock-smart"
          value={form.model}
          onChange={set("model")}
          error={errors.model}
        />

        <Textarea
          label="Системный промпт"
          placeholder="Инструкция для ассистента..."
          value={form.systemPrompt}
          rows={5}
          onChange={set("systemPrompt")}
          error={errors.systemPrompt}
        />

        <Textarea
          label="Пример запроса пользователя"
          placeholder="Необязательно"
          value={form.exampleUserPrompt}
          rows={2}
          onChange={set("exampleUserPrompt")}
        />

        <label className="new-assistant__checkbox">
          <input
            type="checkbox"
            checked={form.isActive}
            onChange={set("isActive")}
          />
          Активен
        </label>

        {serverError && (
          <p className="new-assistant__error">{serverError}</p>
        )}

        <div className="new-assistant__actions">
          <Link to="/assistants">
            <Button variant="secondary">Отмена</Button>
          </Link>
          <Button loading={loading} onClick={handleSubmit}>
            Создать
          </Button>
        </div>
      </div>
    </div>
  );
}