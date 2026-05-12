import axios from "axios";
import useAuthStore from "@/store/auth";

const client = axios.create({
  baseURL: "/api",
  headers: {
    "Content-Type": "application/json",
  },
});

client.interceptors.request.use((config) => {
  const token = useAuthStore.getState().token;
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

client.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      useAuthStore.getState().clearAuth();
      const reqUrl = String(error.config?.url || "");
      const isCredentialEndpoint =
        reqUrl.endsWith("/login") ||
        reqUrl.endsWith("/register") ||
        reqUrl.endsWith("/dummyLogin");
      if (!isCredentialEndpoint) {
        window.location.href = "/login";
      }
    }
    return Promise.reject(error);
  }
);

export default client;