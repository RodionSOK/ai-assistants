import '../Input/Input.css';
import './Select.css';
import clsx from "clsx";

export default function Select({
  label,
  error,
  hint,
  options = [],
  placeholder,
  className,
  ...props
}) {
  return (
    <div className="field">
      {label && <label className="field-label">{label}</label>}
      <select
        className={clsx("field-input field-select", error && "field-input--error", className)}
        {...props}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((opt) => (
          <option key={opt.value} value={opt.value}>
            {opt.label}
          </option>
        ))}
      </select>
      {error && <span className="field-error">{error}</span>}
      {hint && !error && <span className="field-hint">{hint}</span>}
    </div>
  );
}