import { Navigate, Outlet } from "react-router-dom";
import useAuthStore from "@/store/auth";
import Navbar from "@/components/Navbar";

export default function PrivateRoute() {
  const token = useAuthStore((s) => s.token);

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return (
    <>
      <Navbar />
      <main>
        <Outlet />
      </main>
    </>
  );
}