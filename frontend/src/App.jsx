import Router from "@/router";
import useThemeStore from "@/store/theme";
import { useEffect } from "react";

export default function App() {
  const theme = useThemeStore((state) => state.theme);

  useEffect(() => {
    document.documentElement.setAttribute("data-theme", theme);
  }, [theme]);

  return <Router />;
}