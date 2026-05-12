import { useState } from "react";
import { useNavigate, Link } from "react-router-dom";
import { createCategory } from "@/api/categories";
import Input from "@/components/ui/Input";
import Textarea from "@/components/ui/Textarea";
import Button from "@/components/ui/Button";
import "./NewCategory.css";

export default function NewCategory() {
  const navigate = useNavigate();

  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const [errors, setErrors] = useState({});
  const [loading, setLoading] = useState(false);
  const [serverError, setServerError] = useState(null);

  const validate = () => {
    const e = {};
    if (!name.trim()) e.name = "Название обязательно";
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
      await createCategory({ name: name.trim(), description: description.trim() || null });
      navigate("/assistants");
    } catch {
      setServerError("Не удалось создать категорию");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="new-category">
      <div className="new-category__back">
        <Link to="/assistants">← Назад</Link>
      </div>

      <div className="new-category__card">
        <h1 className="new-category__title">Новая категория</h1>

        <Input
          placeholder="Название"
          value={name}
          onChange={(e) => setName(e.target.value)}
          error={errors.name}
        />

        <Textarea
          label="Описание"
          placeholder="Краткое описание категории (необязательно)"
          value={description}
          rows={3}
          onChange={(e) => setDescription(e.target.value)}
        />

        {serverError && (
          <p className="new-category__error">{serverError}</p>
        )}

        <div className="new-category__actions">
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