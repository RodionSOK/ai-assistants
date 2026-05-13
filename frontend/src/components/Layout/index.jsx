import { useNavigate } from "react-router-dom";
import Sidebar from "@/components/Sidebar";
import Button from "@/components/ui/Button";
import MoonIcon from "@/assets/icons/MoonIcon";
import SunIcon from "@/assets/icons/SunIcon";
import useThemeStore from "@/store/theme";
import "./Layout.css";

export default function Layout({ children }) {
  const { theme, toggleTheme } = useThemeStore();

  return (
    <div className="layout">
      <Sidebar />
      <Button
        variant="iconGlass"
        className="layout__theme-btn"
        onClick={toggleTheme}
        title={theme === "light" ? "Тёмная тема" : "Светлая тема"}
        aria-label={theme === "light" ? "Включить тёмную тему" : "Включить светлую тему"}
      >
        {theme === "light" ? <MoonIcon size={20} /> : <SunIcon size={20} />}
      </Button>
      <main className="layout__main">
        {children}
      </main>
    </div>
  );
}