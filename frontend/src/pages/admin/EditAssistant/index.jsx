import { useState, useEffect } from "react";
import { useNavigate, useParams, Link } from "react-router-dom";
import { getAssistant, updateAssistant } from "@/api/assistants";
import { getCategories } from "@/api/categories";
import Input from "@/components/ui/Input";
import Textarea from "@/components/ui/Textarea";
import Select from "@/components/ui/Select";
import Button from "@/components/ui/Button";
import Spinner from "@/components/ui/Spinner";
import "./EditAssistant.css";

export default function EditAssistant() {
  const { id } = useParams();
  const navigate = useNavigate();

  const [categories, setCategories] = useState([]);
  const [form, setForm] = useState(null);
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [serverError, setServerError] = useState(null);

  useEffect(() => {
    Promise.all([getAssistant(id), getCategories()])
      .then(([assistant, { categories }]) => {
        setCategories(categories);
        setForm({
          categoryId: assistant.categoryId,
          name: assistant.name,
          description: assistant.description,
          model: assistant.model,
          systemPrompt: assistant.systemPrompt || "",
          exampleUserPrompt: assistant.exampleUserPrompt || "",
          isActive: assistant.isActive,
        });
      })
      .catch(() => setServerError("Не удалось загрузить ассистента"))
      .finally(() => setLoading(false));
  }, [id]);

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

    setSaving(true);
    setServerError(null);

    try {
      await updateAssistant(id, {
        ...form,
        name: form.name.trim(),
        description: form.description.trim(),
        model: form.model.trim(),
        systemPrompt: form.systemPrompt.trim(),
        exampleUserPrompt: form.exampleUserPrompt.trim() || null,
      });
      navigate(`/assistants/${id}`);
    } catch (e) {
      const msg = e.response?.data?.error?.message;
      setServerError(msg || "Не удалось сохранить изменения");
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="edit-assistant__center">
        <Spinner size="lg" />
      </div>
    );
  }

  if (!form) {
    return (
      <div className="edit-assistant__center">
        <p className="edit-assistant__error">{serverError}</p>
        <Link to="/assistants">← Назад к каталогу</Link>
      </div>
    );
  }

  return (
    <div className="edit-assistant">
      <div className="edit-assistant__back">
        <Link to={`/assistants/${id}`}>← Назад</Link>
      </div>

      <div className="edit-assistant__card">
        <h1 className="edit-assistant__title">Редактировать ассистента</h1>

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
          value={form.name}
          onChange={set("name")}
          error={errors.name}
        />

        <Textarea
          label="Описание"
          value={form.description}
          rows={2}
          onChange={set("description")}
          error={errors.description}
        />

        <Input
          label="Модель"
          value={form.model}
          onChange={set("model")}
          error={errors.model}
        />

        <Textarea
          label="Системный промпт"
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

        <label className="edit-assistant__checkbox">
          <input
            type="checkbox"
            checked={form.isActive}
            onChange={set("isActive")}
          />
          Активен
        </label>

        {serverError && (
          <p className="edit-assistant__error">{serverError}</p>
        )}

        <div className="edit-assistant__actions">
          <Link to={`/assistants/${id}`}>
            <Button variant="secondary">Отмена</Button>
          </Link>
          <Button loading={saving} onClick={handleSubmit}>
            Сохранить
          </Button>
        </div>
      </div>
    </div>
  );
}