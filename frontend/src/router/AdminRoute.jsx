import { Navigate, Outlet } from "react-router-dom";
import useAuthStore from "@/store/auth";

export default function AdminRoute() {
  const user = useAuthStore((state) => state.user);

  if (user?.role !== "admin") {
    return <Navigate to="/assistants" replace />;
  }

  return <Outlet />;
}