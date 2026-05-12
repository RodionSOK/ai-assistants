import { Routes, Route, Navigate } from "react-router-dom";
import useAuthStore from "@/store/auth";

import Login from "@/pages/Login";
import Assistants from "@/pages/Assistants";
import AssistantDetail from "@/pages/AssistantDetail";
import MyRuns from "@/pages/MyRuns";
import NewCategory from "@/pages/admin/NewCategory";
import NewAssistant from "@/pages/admin/NewAssistant";
import EditAssistant from "@/pages/admin/EditAssistant";
import AllRuns from "@/pages/admin/AllRuns";

import PrivateRoute from "./PrivateRoute";
import AdminRoute from "./AdminRoute";

export default function Router() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />

      <Route element={<PrivateRoute />}>
        <Route path="/assistants" element={<Assistants />} />
        <Route path="/assistants/:id" element={<AssistantDetail />} />
        <Route path="/runs/my" element={<MyRuns />} />

        <Route element={<AdminRoute />}>
          <Route path="/admin/categories/new" element={<NewCategory />} />
          <Route path="/admin/assistants/new" element={<NewAssistant />} />
          <Route path="/admin/assistants/:id/edit" element={<EditAssistant />} />
          <Route path="/admin/runs" element={<AllRuns />} />
        </Route>
      </Route>

      <Route path="/" element={<Navigate to="/assistants" replace />} />
      <Route path="*" element={<Navigate to="/assistants" replace />} />
    </Routes>
  );
}