import './Empty.css';
export default function Empty({ title = "Ничего не найдено", description }) {
  return (
    <div className="empty">
      <div className="empty-icon">🔍</div>
      <p className="empty-title">{title}</p>
      {description && <p className="empty-description">{description}</p>}
    </div>
  );
}