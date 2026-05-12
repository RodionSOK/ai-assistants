import '../Input/Input.css';
import './Textarea.css';
import clsx from "clsx";

export default function Textarea({
  label,
  error,
  hint,
  rows = 4,
  className,
  ...props
}) {
  return (
    <div className="field">
      {label && <label className="field-label">{label}</label>}
      <textarea
        rows={rows}
        className={clsx("field-input field-textarea", error && "field-input--error", className)}
        {...props}
      />
      {error && <span className="field-error">{error}</span>}
      {hint && !error && <span className="field-hint">{hint}</span>}
    </div>
  );
}