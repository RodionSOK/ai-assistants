import './Input.css';
import clsx from "clsx";

export default function Input({
  label,
  error,
  hint,
  className,
  ...props
}) {
  return (
    <div className="field">
      {label && <label className="field-label">{label}</label>}
      <input
        className={clsx("field-input", error && "field-input--error", className)}
        {...props}
      />
      {error && <span className="field-error">{error}</span>}
      {hint && !error && <span className="field-hint">{hint}</span>}
    </div>
  );
}