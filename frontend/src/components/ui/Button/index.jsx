import "./Button.css";
import clsx from "clsx";

const variants = {
  primary: "btn-primary",
  secondary: "btn-secondary",
  ghost: "btn-ghost",
  segment: "btn-segment",
  iconGlass: "btn-icon-glass",
  role: "btn-role",
};

const sizes = {
  sm: "btn-sm",
  md: "btn-md",
  lg: "btn-lg",
};

const layoutVariants = new Set(["segment", "iconGlass", "role"]);

export default function Button({
  children,
  variant = "primary",
  size = "md",
  pressed = false,
  disabled = false,
  loading = false,
  type = "button",
  onClick,
  className,
  ...props
}) {
  const sizeClass = layoutVariants.has(variant) ? null : sizes[size];

  return (
    <button
      type={type}
      disabled={disabled || loading}
      onClick={onClick}
      className={clsx(
        "btn",
        variants[variant],
        variant === "segment" && pressed && "btn-segment--active",
        variant === "role" && pressed && "btn-role--active",
        sizeClass,
        className
      )}
      {...props}
    >
      {loading ? <span className="btn-spinner" /> : children}
    </button>
  );
}