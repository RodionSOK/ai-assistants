import { render, screen } from "@testing-library/react";
import StatusBadge from "./index";

describe("StatusBadge", () => {
  test("показывает 'Успешно' для статуса success", () => {
    render(<StatusBadge status="success" />);
    const badge = screen.getByText("Успешно");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("badge--success");
  });

  test("показывает 'Ошибка' для статуса failed", () => {
    render(<StatusBadge status="failed" />);
    const badge = screen.getByText("Ошибка");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("badge--danger");
  });

  test("показывает 'В обработке' для статуса pending", () => {
    render(<StatusBadge status="pending" />);
    const badge = screen.getByText("В обработке");
    expect(badge).toBeInTheDocument();
    expect(badge).toHaveClass("badge--warning");
  });

  test("показывает сырой статус для неизвестного значения", () => {
    render(<StatusBadge status="unknown_status" />);
    expect(screen.getByText("unknown_status")).toBeInTheDocument();
  });
});