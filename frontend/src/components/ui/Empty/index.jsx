import './Empty.css';
export default function Empty({ title = "Ничего не найдено", description }) {
  return (
    <div className="empty">
      <p className="empty-title">{title}</p>
      {description && <p className="empty-description">{description}</p>}
    </div>
  );
}