import { useState, useEffect } from "react";
import { useParams, useNavigate, Link } from "react-router-dom";
import { getAssistant } from "@/api/assistants";
import { runAssistant } from "@/api/runs";
import useAuthStore from "@/store/auth";
import StatusBadge from "@/components/StatusBadge";
import Button from "@/components/ui/Button";
import Spinner from "@/components/ui/Spinner";
import Textarea from "@/components/ui/Textarea";
import "./AssistantDetail.css";

export default function AssistantDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const user = useAuthStore((s) => s.user);
  const isAdmin = user?.role === "admin";

  const [assistant, setAssistant] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [prompt, setPrompt] = useState("");
  const [running, setRunning] = useState(false);
  const [runResult, setRunResult] = useState(null);
  const [runError, setRunError] = useState(null);

  useEffect(() => {
    setLoading(true);
    getAssistant(id)
      .then((data) => setAssistant(data))
      .catch(() => setError("Не удалось загрузить ассистента"))
      .finally(() => setLoading(false));
  }, [id]);

  const handleRun = async () => {
    if (!prompt.trim()) return;
    setRunning(true);
    setRunResult(null);
    setRunError(null);
    try {
      const result = await runAssistant(id, prompt.trim());
      setRunResult(result);
    } catch (e) {
      setRunError(e?.response?.data?.error || "Ошибка при запуске ассистента");
    } finally {
      setRunning(false);
    }
  };

  if (loading) {
    return (
      <div className="assistant-detail">
        <div className="assistant-detail__center">
          <Spinner size="lg" />
        </div>
      </div>
    );
  }

  if (error || !assistant) {
    return (
      <div className="assistant-detail">
        <div className="assistant-detail__center">
          <p>{error || "Ассистент не найден"}</p>
          <Button variant="secondary" onClick={() => navigate("/assistants")}>
            Вернуться к каталогу
          </Button>
        </div>
      </div>
    );
  }

  return (
    <div className="assistant-detail">
      <Link to="/assistants" className="assistant-detail__back">
        ← Назад к каталогу
      </Link>

      <div className="assistant-detail__card">
        <div className="assistant-detail__header">
          <div>
            <span className="assistant-detail__category">
              {assistant.categoryName}
            </span>
            <h1 className="assistant-detail__name">{assistant.name}</h1>
          </div>
          <div className="assistant-detail__meta">
            {assistant.isActive
                ? <span className="badge badge--success">Активен</span>
                : <span className="badge badge--default">Неактивен</span>}
            {isAdmin && (
              <Button
                variant="secondary"
                size="sm"
                onClick={() => navigate(`/admin/assistants/${id}/edit`)}
              >
                Редактировать
              </Button>
            )}
          </div>
        </div>

        {assistant.description && (
          <p className="assistant-detail__description">
            {assistant.description}
          </p>
        )}

        <div className="assistant-detail__info">
          <span className="assistant-detail__model">
            Модель: {assistant.model}
          </span>
        </div>
      </div>

      <div className="assistant-detail__run">
        <h2 className="assistant-detail__run-title">Запустить ассистента</h2>

        {!assistant.isActive && (
          <p className="assistant-detail__inactive-warn">
            Ассистент неактивен и недоступен для запуска.
          </p>
        )}

        <Textarea
          placeholder="Введите ваш запрос..."
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          rows={4}
          disabled={!assistant.isActive || running}
        />

        <Button
          onClick={handleRun}
          loading={running}
          disabled={!assistant.isActive || !prompt.trim()}
        >
          Запустить
        </Button>

        {runError && (
          <p className="assistant-detail__run-error">{runError}</p>
        )}

        {runResult && (
          <div className="assistant-detail__result">
            <div className="assistant-detail__result-header">
              <span className="assistant-detail__result-label">Результат</span>
              <StatusBadge status={runResult.status} />
            </div>
            {runResult.output && (
              <pre className="assistant-detail__output">{runResult.output}</pre>
            )}
            {runResult.error && (
              <p className="assistant-detail__error">{runResult.error}</p>
            )}
          </div>
        )}
      </div>
    </div>
  );
}