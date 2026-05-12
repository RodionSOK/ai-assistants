import clsx from "clsx";
import "./Icon.css";

const svgRawByName = import.meta.glob("../../assets/icons/*.svg", {
  eager: true,
  query: "?raw",
  import: "default",
});

function resolveSvgString(name) {
  const safe = String(name).replace(/[^a-z0-9_-]/gi, "");
  if (!safe) return null;
  const suffix = `/${safe}.svg`;
  const entry = Object.entries(svgRawByName).find(([path]) => path.endsWith(suffix));
  return entry ? entry[1] : null;
}

export default function Icon({
  name,
  size = 20,
  color,
  className,
  title,
  decorative = true,
}) {
  const html = resolveSvgString(name);
  if (!html) {
    if (import.meta.env.DEV) {
      console.warn(`[Icon] Нет файла assets/icons/${name}.svg`);
    }
    return null;
  }

  return (
    <span
      className={clsx("icon", className)}
      style={{
        width: size,
        height: size,
        ...(color != null && color !== "" ? { color } : {}),
      }}
      role={decorative ? undefined : "img"}
      aria-hidden={decorative ? true : undefined}
      aria-label={decorative ? undefined : title}
      title={title}
      dangerouslySetInnerHTML={{ __html: html }}
    />
  );
}
