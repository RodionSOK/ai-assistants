import { create } from "zustand";

const THEME_KEY = "ai_assistants_theme";

const useThemeStore = create((set) => ({
  theme: localStorage.getItem(THEME_KEY) || "light",

  toggleTheme: () =>
    set((state) => {
      const next = state.theme === "light" ? "dark" : "light";
      localStorage.setItem(THEME_KEY, next);
      return { theme: next };
    }),
}));

export default useThemeStore;