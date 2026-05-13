import { Navigate, Outlet } from "react-router-dom";
import useAuthStore from "@/store/auth";
import Layout from "@/components/Layout";

export default function PrivateRoute() {
  const token = useAuthStore((state) => state.token);

  if (!token) {
    return <Navigate to="/login" replace />;
  }

  return (
    <Layout>
      <Outlet />
    </Layout>
  );
}