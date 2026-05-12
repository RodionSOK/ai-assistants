import { Link, useNavigate, useLocation } from "react-router-dom";
import useAuthStore from "@/store/auth";
import useThemeStore from "@/store/theme";
import Button from "@/components/ui/Button";
import "./Navbar.css";

export default function Navbar() {
  const user = useAuthStore((s) => s.user);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const { theme, toggleTheme } = useThemeStore();
  const navigate = useNavigate();
  const location = useLocation();

  const isAdmin = user?.role === "admin";

  const isActive = (path) => location.pathname.startsWith(path);

  const handleLogout = () => {
    clearAuth();
    navigate("/login", { replace: true });
  };

  return (
    <header className="navbar">
      <div className="navbar__inner">
        <div className="navbar__left">
          <Link to="/assistants" className="navbar__logo">
            AI Ассистенты
          </Link>
          <nav className="navbar__nav">
            <Link
              to="/assistants"
              className={`navbar__link ${isActive("/assistants") ? "navbar__link--active" : ""}`}
            >
              Каталог
            </Link>
            <Link
              to="/runs/my"
              className={`navbar__link ${isActive("/runs") ? "navbar__link--active" : ""}`}
            >
              Мои запуски
            </Link>
            {isAdmin && (
              <>
                <Link
                  to="/admin/runs"
                  className={`navbar__link ${isActive("/admin/runs") ? "navbar__link--active" : ""}`}
                >
                  Все запуски
                </Link>
                <Link
                  to="/admin/categories/new"
                  className={`navbar__link ${isActive("/admin/categories") ? "navbar__link--active" : ""}`}
                >
                  + Категория
                </Link>
                <Link
                  to="/admin/assistants/new"
                  className={`navbar__link ${isActive("/admin/assistants") ? "navbar__link--active" : ""}`}
                >
                  + Ассистент
                </Link>
              </>
            )}
          </nav>
        </div>

        <div className="navbar__right">
          <button className="navbar__theme-btn" onClick={toggleTheme} title="Переключить тему">
            {theme === "light" ? "🌙" : "☀️"}
          </button>
          <span className="navbar__user">{user?.email}</span>
          <Button variant="secondary" size="sm" onClick={handleLogout}>
            Выйти
          </Button>
        </div>
      </div>
    </header>
  );
}