import './AssistantCard.css';
import { Link } from "react-router-dom";
import clsx from "clsx";

export default function AssistantCard({ assistant }) {
  return (
    <Link
      to={`/assistants/${assistant.id}`}
      className={clsx("assistant-card", !assistant.isActive && "assistant-card--inactive")}
    >
      <div className="assistant-card__header">
        <span className="assistant-card__category">{assistant.categoryName}</span>
        {!assistant.isActive && (
          <span className="badge badge--default">Неактивен</span>
        )}
      </div>
      <h3 className="assistant-card__name">{assistant.name}</h3>
      <p className="assistant-card__description">{assistant.description}</p>
      <div className="assistant-card__footer">
        <span className="assistant-card__model">{assistant.model}</span>
      </div>
    </Link>
  );
}