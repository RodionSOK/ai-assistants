import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { dummyLogin, register, login } from "@/api/auth";
import useAuthStore from "@/store/auth";
import useThemeStore from "@/store/theme";
import MoonIcon from "@/assets/icons/MoonIcon";
import SunIcon from "@/assets/icons/SunIcon";
import AdminIcon from "@/assets/icons/AdminIcon";
import UserIcon from "@/assets/icons/UserIcon";
import Button from "@/components/ui/Button";
import Input from "@/components/ui/Input";
import "./Login.css";

const TABS = [
  { id: "login", label: "Войти" },
  { id: "register", label: "Регистрация" },
  { id: "demo", label: "Demo" },
];

export default function Login() {
  const [tab, setTab] = useState("login");
  const { theme, toggleTheme } = useThemeStore();

  const setAuth = useAuthStore((s) => s.setAuth);
  const navigate = useNavigate();

  const handleSuccess = (token, user) => {
    setAuth(token, user);
    navigate("/assistants", { replace: true });
  };

  return (
    <div className="login">
      <Button
        variant="iconGlass"
        className="login__theme"
        onClick={toggleTheme}
        title={theme === "light" ? "Тёмная тема" : "Светлая тема"}
        aria-label={theme === "light" ? "Включить тёмную тему" : "Включить светлую тему"}
      >
        {theme === "light" ? <MoonIcon size={20} /> : <SunIcon size={20} />}
      </Button>

      <div className="login__card">
        <header className="login__header">
          <h1 className="login__title">AI Ассистенты</h1>
        </header>

        <div className="login__tabs" role="tablist" aria-label="Способ входа">
          {TABS.map((t) => (
            <Button
              key={t.id}
              variant="segment"
              pressed={tab === t.id}
              role="tab"
              aria-selected={tab === t.id}
              onClick={() => setTab(t.id)}
            >
              {t.label}
            </Button>
          ))}
        </div>

        {tab === "login" && <LoginForm onSuccess={handleSuccess} />}
        {tab === "register" && <RegisterForm onSuccess={handleSuccess} />}
        {tab === "demo" && <DemoForm onSuccess={handleSuccess} />}
      </div>
    </div>
  );
}

function LoginForm({ onSuccess }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState(null);
  const [loading, setLoading] = useState(false);

  const validate = () => {
    const e = {};
    if (!email) e.email = "Введите email";
    else if (!/\S+@\S+\.\S+/.test(email)) e.email = "Некорректный email";
    if (!password) e.password = "Введите пароль";
    return e;
  };

  const handleSubmit = async () => {
    const e = validate();
    if (Object.keys(e).length) { setErrors(e); return; }

    setLoading(true);
    setServerError(null);
    setErrors({});
    try {
      const { token, user } = await login(email, password);
      onSuccess(token, user);
    } catch (err) {
      const msg = err.response?.data?.error?.message;
      setServerError(msg || "Неверный email или пароль");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login__form">
      <Input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        error={errors.email}
        autoComplete="email"
      />
      <Input
        type="password"
        placeholder="Пароль"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        error={errors.password}
        autoComplete="current-password"
        onKeyDown={(e) => e.key === "Enter" && handleSubmit()}
      />
      {serverError && <p className="login__error">{serverError}</p>}
      <Button size="lg" loading={loading} onClick={handleSubmit} className="login__btn">
        Войти
      </Button>
    </div>
  );
}

function RegisterForm({ onSuccess }) {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirm, setConfirm] = useState("");
  const [errors, setErrors] = useState({});
  const [serverError, setServerError] = useState(null);
  const [loading, setLoading] = useState(false);

  const validate = () => {
    const e = {};
    if (!email) e.email = "Введите email";
    else if (!/\S+@\S+\.\S+/.test(email)) e.email = "Некорректный email";
    if (!password) e.password = "Введите пароль";
    else if (password.length < 8) e.password = "Минимум 8 символов";
    if (!confirm) e.confirm = "Подтвердите пароль";
    else if (confirm !== password) e.confirm = "Пароли не совпадают";
    return e;
  };

  const handleSubmit = async () => {
    const e = validate();
    if (Object.keys(e).length) { setErrors(e); return; }

    setLoading(true);
    setServerError(null);
    setErrors({});
    try {
      const { token, user } = await register(email, password);
      onSuccess(token, user);
    } catch (err) {
      const msg = err.response?.data?.error?.message;
      setServerError(msg || "Не удалось зарегистрироваться. Попробуйте ещё раз.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login__form">
      <Input
        type="email"
        placeholder="Email"
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        error={errors.email}
        autoComplete="email"
      />
      <Input
        type="password"
        placeholder="Пароль"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        error={errors.password}
        autoComplete="new-password"
      />
      <Input
        type="password"
        placeholder="Подтверждение пароля"
        value={confirm}
        onChange={(e) => setConfirm(e.target.value)}
        error={errors.confirm}
        autoComplete="new-password"
        onKeyDown={(e) => e.key === "Enter" && handleSubmit()}
      />
      {serverError && <p className="login__error">{serverError}</p>}
      <Button size="lg" loading={loading} onClick={handleSubmit} className="login__btn">
        Зарегистрироваться
      </Button>
    </div>
  );
}

function DemoForm({ onSuccess }) {
  const [role, setRole] = useState("user");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleLogin = async () => {
    setLoading(true);
    setError(null);
    try {
      const { token, user } = await dummyLogin(role);
      onSuccess(token, user);
    } catch {
      setError("Не удалось войти. Попробуйте ещё раз.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="login__form">
      <div className="login__roles login__roles--icon-only">
        <Button
          variant="role"
          className="login__role-icon-only"
          pressed={role === "user"}
          onClick={() => setRole("user")}
          aria-label="Режим пользователя"
        >
          <UserIcon size={40} />
        </Button>
        <Button
          variant="role"
          className="login__role-icon-only"
          pressed={role === "admin"}
          onClick={() => setRole("admin")}
          aria-label="Режим администратора"
        >
          <AdminIcon size={40} />
        </Button>
      </div>
      {error && <p className="login__error">{error}</p>}
      <Button
        size="lg"
        loading={loading}
        onClick={handleLogin}
        className="login__btn"
        aria-label={
          role === "admin"
            ? "Войти в демо как администратор"
            : "Войти в демо как пользователь"
        }
      >
        {role === "admin" ? "Войти как администратор" : "Войти как пользователь"}
      </Button>
    </div>
  );
}