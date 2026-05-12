import './StatusBadge.css';
const config = {
  pending: { label: "В обработке", className: "badge badge--warning" },
  success: { label: "Успешно", className: "badge badge--success" },
  failed: { label: "Ошибка", className: "badge badge--danger" },
};

export default function StatusBadge({ status }) {
  const { label, className } = config[status] || {
    label: status,
    className: "badge badge--default",
  };

  return <span className={className}>{label}</span>;
}