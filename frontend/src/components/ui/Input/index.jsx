import "./Input.css";
import clsx from "clsx";

export default function Input({
  error,
  hint,
  className,
  placeholder,
  "aria-label": ariaLabel,
  ...props
}) {
  const a11yLabel = ariaLabel ?? placeholder;

  return (
    <div className="field">
      <input
        placeholder={placeholder}
        aria-label={a11yLabel}
        aria-invalid={error ? true : undefined}
        className={clsx("field-input", error && "field-input--error", className)}
        {...props}
      />
      {error && <span className="field-error">{error}</span>}
      {hint && !error && <span className="field-hint">{hint}</span>}
    </div>
  );
}
