import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { getAssistants } from "@/api/assistants";
import { getCategories } from "@/api/categories";
import useAuthStore from "@/store/auth";
import AssistantCard from "@/components/AssistantCard";
import Pagination from "@/components/Pagination";
import Spinner from "@/components/ui/Spinner";
import Empty from "@/components/ui/Empty";
import "./Assistants.css";

export default function Assistants() {
  const [searchParams, setSearchParams] = useSearchParams();
  const user = useAuthStore((s) => s.user);
  const isAdmin = user?.role === "admin";

  const [assistants, setAssistants] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, pageSize: 12, total: 0 });
  const [categories, setCategories] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const page = parseInt(searchParams.get("page") || "1");
  const categoryId = searchParams.get("categoryId") || "";
  const q = searchParams.get("q") || "";
  const includeInactive = searchParams.get("includeInactive") === "true";

  useEffect(() => {
    getCategories()
      .then(({ categories }) => setCategories(categories))
      .catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    setError(null);

    const params = { page, pageSize: 12 };
    if (categoryId) params.categoryId = categoryId;
    if (q) params.q = q;
    if (isAdmin && includeInactive) params.includeInactive = true;

    getAssistants(params)
      .then(({ assistants, pagination }) => {
        setAssistants(assistants);
        setPagination(pagination);
      })
      .catch(() => setError("Не удалось загрузить ассистентов"))
      .finally(() => setLoading(false));
  }, [page, categoryId, q, includeInactive]);

  const setParam = (key, value) => {
    setSearchParams((prev) => {
      const next = new URLSearchParams(prev);
      if (value) {
        next.set(key, value);
      } else {
        next.delete(key);
      }
      next.delete("page");
      return next;
    });
  };

  const handleSearch = (e) => {
    setParam("q", e.target.value);
  };

  return (
    <div className="assistants">
      <div className="assistants__header">
        <h1 className="assistants__title">Каталог ассистентов</h1>
        {isAdmin && (
          <label className="assistants__inactive-toggle">
            <input
              type="checkbox"
              checked={includeInactive}
              onChange={(e) => setParam("includeInactive", e.target.checked ? "true" : "")}
            />
            Показать неактивных
          </label>
        )}
      </div>

      <div className="assistants__filters">
        <input
          className="assistants__search"
          type="search"
          placeholder="Поиск по названию и описанию..."
          value={q}
          onChange={handleSearch}
        />
        <div className="assistants__categories">
          <button
            className={`assistants__cat-btn ${!categoryId ? "assistants__cat-btn--active" : ""}`}
            onClick={() => setParam("categoryId", "")}
          >
            Все
          </button>
          {categories.map((cat) => (
            <button
              key={cat.id}
              className={`assistants__cat-btn ${categoryId === cat.id ? "assistants__cat-btn--active" : ""}`}
              onClick={() => setParam("categoryId", cat.id)}
            >
              {cat.name}
            </button>
          ))}
        </div>
      </div>

      {loading && (
        <div className="assistants__spinner">
          <Spinner size="lg" />
        </div>
      )}

      {error && !loading && (
        <p className="assistants__error">{error}</p>
      )}

      {!loading && !error && assistants.length === 0 && (
        <Empty
          title="Ассистенты не найдены"
          description="Попробуйте изменить фильтры или поисковый запрос"
        />
      )}

      {!loading && !error && assistants.length > 0 && (
        <>
          <div className="assistants__grid">
            {assistants.map((a) => (
              <AssistantCard key={a.id} assistant={a} />
            ))}
          </div>
          <Pagination
            page={pagination.page}
            pageSize={pagination.pageSize}
            total={pagination.total}
            onPageChange={(p) =>
              setSearchParams((prev) => {
                const next = new URLSearchParams(prev);
                next.set("page", p);
                return next;
              })
            }
          />
        </>
      )}
    </div>
  );
}