import { create } from "zustand";

const TOKEN_KEY = "ai_assistants_token";
const USER_KEY = "ai_assistants_user";

const useAuthStore = create((set) => ({
  token: localStorage.getItem(TOKEN_KEY) || null,
  user: JSON.parse(localStorage.getItem(USER_KEY) || "null"),

  setAuth: (token, user) => {
    localStorage.setItem(TOKEN_KEY, token);
    localStorage.setItem(USER_KEY, JSON.stringify(user));
    set({ token, user });
  },

  clearAuth: () => {
    localStorage.removeItem(TOKEN_KEY);
    localStorage.removeItem(USER_KEY);
    set({ token: null, user: null });
  },

  isAuthenticated: () => {
    return !!localStorage.getItem(TOKEN_KEY);
  },
}));

export default useAuthStore;