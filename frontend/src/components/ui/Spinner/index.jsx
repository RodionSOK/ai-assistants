import './Spinner.css';
import clsx from "clsx";

export default function Spinner({ size = "md", className }) {
  return (
    <div className={clsx("spinner", `spinner--${size}`, className)} />
  );
}