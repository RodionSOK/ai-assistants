import { Link, useNavigate, useLocation } from "react-router-dom";
import useAuthStore from "@/store/auth";
import Button from "@/components/ui/Button";
import "./Sidebar.css";

const NAV_ITEMS = [
  { path: "/assistants", label: "Каталог", exact: false },
  { path: "/runs/my", label: "Мои запуски", exact: false },
];

const ADMIN_NAV_ITEMS = [
  { path: "/admin/runs", label: "Все запуски", exact: false },
  { path: "/admin/categories/new", label: "Новая категория", exact: false },
  { path: "/admin/assistants/new", label: "Новый ассистент", exact: false },
];

export default function Sidebar() {
  const user = useAuthStore((s) => s.user);
  const clearAuth = useAuthStore((s) => s.clearAuth);
  const navigate = useNavigate();
  const location = useLocation();

  const isAdmin = user?.role === "admin";

  const isActive = (path) => location.pathname.startsWith(path);

  const handleLogout = () => {
    clearAuth();
    navigate("/login", { replace: true });
  };

  return (
    <aside className="sidebar">
      <div className="sidebar__top">
        <Link to="/assistants" className="sidebar__logo">
          AI Ассистенты
        </Link>

        <nav className="sidebar__nav">
          <p className="sidebar__nav-label">Основное</p>
          {NAV_ITEMS.map((item) => (
            <Link
              key={item.path}
              to={item.path}
              className={`sidebar__link ${isActive(item.path) ? "sidebar__link--active" : ""}`}
            >
              {item.label}
            </Link>
          ))}

          {isAdmin && (
            <>
              <p className="sidebar__nav-label">Администратор</p>
              {ADMIN_NAV_ITEMS.map((item) => (
                <Link
                  key={item.path}
                  to={item.path}
                  className={`sidebar__link ${isActive(item.path) ? "sidebar__link--active" : ""}`}
                >
                  {item.label}
                </Link>
              ))}
            </>
          )}
        </nav>
      </div>

      <div className="sidebar__bottom">
        <div className="sidebar__user">
            <span className="sidebar__user-role">
            {isAdmin ? "Администратор" : "Пользователь"}
            </span>
            <span className="sidebar__user-email">{user?.email}</span>
        </div>
        <div className="sidebar__actions">
            <Button size="sm" onClick={handleLogout} className="logout__btn" >
                Выйти
            </Button>
        </div>
        </div>
    </aside>
  );
}