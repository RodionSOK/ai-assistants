import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { dummyLogin } from "@/api/auth";
import useAuthStore from "@/store/auth";
import Button from "@/components/ui/Button";
import "./Login.css";

export default function Login() {
  const [role, setRole] = useState("user");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const setAuth = useAuthStore((s) => s.setAuth);
  const navigate = useNavigate();

  const handleLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      const { token, user } = await dummyLogin(role);
      setAuth(token, user);
      navigate("/assistants", { replace: true });
    } catch {
      setError("Не удалось войти. Попробуйте ещё раз.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login">
      <div className="login__card">
        <h1 className="login__title">AI Ассистенты</h1>
        <p className="login__subtitle">Выберите роль для входа</p>

        <div className="login__roles">
          <button
            className={`login__role ${role === "user" ? "login__role--active" : ""}`}
            onClick={() => setRole("user")}
          >
            <span className="login__role-icon">👤</span>
            <span className="login__role-label">Пользователь</span>
          </button>
          <button
            className={`login__role ${role === "admin" ? "login__role--active" : ""}`}
            onClick={() => setRole("admin")}
          >
            <span className="login__role-icon">🛠️</span>
            <span className="login__role-label">Администратор</span>
          </button>
        </div>

        {error && <p className="login__error">{error}</p>}

        <Button size="lg" loading={loading} onClick={handleLogin} className="login__btn">
          Войти как {role === "admin" ? "администратор" : "пользователь"}
        </Button>
      </div>
    </div>
  );
}