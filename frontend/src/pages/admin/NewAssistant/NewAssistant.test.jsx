import { render, screen, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { MemoryRouter } from "react-router-dom";
import { vi } from "vitest";
import NewAssistant from "./index";

vi.mock("@/api/categories", () => ({
  getCategories: () => Promise.resolve({
    categories: [{ id: "cat-1", name: "AI Tools" }],
  }),
}));

vi.mock("@/api/assistants", () => ({
  createAssistant: vi.fn(() => Promise.resolve({})),
}));

const mockNavigate = vi.fn();
vi.mock("react-router-dom", async (importOriginal) => {
  const actual = await importOriginal();
  return { ...actual, useNavigate: () => mockNavigate };
});

const renderPage = () =>
  render(
    <MemoryRouter>
      <NewAssistant />
    </MemoryRouter>
  );

describe("NewAssistant — валидация формы", () => {
  test("показывает ошибки при сабмите пустой формы", async () => {
    renderPage();

    await userEvent.click(screen.getByRole("button", { name: /создать/i }));

    await waitFor(() => {
      expect(screen.getByText("Выберите категорию")).toBeInTheDocument();
      expect(screen.getByText("Название обязательно")).toBeInTheDocument();
      expect(screen.getByText("Описание обязательно")).toBeInTheDocument();
      expect(screen.getByText("Модель обязательна")).toBeInTheDocument();
      expect(screen.getByText("Системный промпт обязателен")).toBeInTheDocument();
    });
  });

  test("ошибка поля сбрасывается при вводе текста", async () => {
    renderPage();

    await userEvent.click(screen.getByRole("button", { name: /создать/i }));
    await waitFor(() => {
      expect(screen.getByText("Название обязательно")).toBeInTheDocument();
    });

    await userEvent.type(screen.getByPlaceholderText("Название"), "Мой ассистент");

    await waitFor(() => {
      expect(screen.queryByText("Название обязательно")).not.toBeInTheDocument();
    });
  });

  test("успешный сабмит вызывает navigate на /assistants", async () => {
    const { createAssistant } = await import("@/api/assistants");
    createAssistant.mockResolvedValueOnce({});

    renderPage();

    await waitFor(() => {
      expect(screen.getByRole("option", { name: "AI Tools" })).toBeInTheDocument();
    });

    await userEvent.selectOptions(
      screen.getByRole("combobox"),
      "cat-1"
    );
    await userEvent.type(screen.getByPlaceholderText("Название"), "Тест");
    await userEvent.type(screen.getByPlaceholderText("Описание"), "Описание");
    await userEvent.type(screen.getByPlaceholderText("Модель"), "gpt-4");
    await userEvent.type(screen.getByPlaceholderText("Системный промпт"), "You are helpful.");

    await userEvent.click(screen.getByRole("button", { name: /создать/i }));

    await waitFor(() => {
      expect(mockNavigate).toHaveBeenCalledWith("/assistants");
    });
  });
});