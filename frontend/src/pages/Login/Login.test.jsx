import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { vi } from "vitest";
import Login from "./index";

vi.mock("@/api/auth", () => ({
  login: vi.fn(),
  register: vi.fn(),
  dummyLogin: vi.fn(),
}));

vi.mock("@/store/auth", () => ({
  default: (selector) =>
    selector({
      setAuth: vi.fn(),
      token: null,
      user: null,
    }),
}));

vi.mock("@/store/theme", () => ({
  default: () => ({ theme: "light", toggleTheme: vi.fn() }),
}));

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async (importOriginal) => {
  const actual = await importOriginal();
  return { ...actual, useNavigate: () => mockNavigate };
});

const renderLogin = () =>
  render(
    <MemoryRouter>
      <Login />
    </MemoryRouter>
  );

describe("Login — LoginForm", () => {
  test("показывает ошибку при пустом email и пароле", async () => {
    renderLogin();

    await userEvent.click(screen.getByRole("button", { name: /войти$/i }));

    await waitFor(() => {
      expect(screen.getByText("Введите email")).toBeInTheDocument();
      expect(screen.getByText("Введите пароль")).toBeInTheDocument();
    });
  });

  test("показывает ошибку при некорректном email", async () => {
    renderLogin();

    await userEvent.type(screen.getByPlaceholderText("Email"), "notanemail");
    await userEvent.click(screen.getByRole("button", { name: /войти$/i }));

    await waitFor(() => {
      expect(screen.getByText("Некорректный email")).toBeInTheDocument();
    });
  });

  test("успешный логин вызывает navigate", async () => {
    const { login } = await import("@/api/auth");
    login.mockResolvedValueOnce({
      token: "test-token",
      user: { id: "1", email: "test@test.com", role: "user" },
    });

    renderLogin();

    await userEvent.type(screen.getByPlaceholderText("Email"), "test@test.com");
    await userEvent.type(screen.getByPlaceholderText("Пароль"), "password123");
    await userEvent.click(screen.getByRole("button", { name: /войти$/i }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith("/assistants", { replace: true });
    });
  });
});