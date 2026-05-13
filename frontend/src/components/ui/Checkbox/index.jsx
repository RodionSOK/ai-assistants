import clsx from "clsx";
import "./Checkbox.css";

export default function Checkbox({ label, checked, onChange, className }) {
  return (
    <label className={clsx("checkbox", className)}>
      <input
        type="checkbox"
        className="checkbox__input"
        checked={checked}
        onChange={onChange}
      />
      <span className="checkbox__box" aria-hidden="true" />
      {label && <span className="checkbox__label">{label}</span>}
    </label>
  );
}