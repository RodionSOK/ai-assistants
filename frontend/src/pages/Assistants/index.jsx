import { useState, useEffect } from "react";
import { useSearchParams, Link } from "react-router-dom";
import { getAssistants } from "@/api/assistants";
import { getCategories } from "@/api/categories";
import useAuthStore from "@/store/auth";
import AssistantCard from "@/components/AssistantCard";
import Pagination from "@/components/Pagination";
import Spinner from "@/components/ui/Spinner";
import Empty from "@/components/ui/Empty";
import Input from "@/components/ui/Input";
import Button from "@/components/ui/Button";
import Checkbox from "@/components/ui/Checkbox";
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

  return (
    <div className="assistants">
      <div className="assistants__header">
        <h1 className="assistants__title">Каталог ассистентов</h1>
      </div>

      <div className="assistants__filters">
        <div className="assistants__search-row">
          <Input
            type="search"
            placeholder="Поиск по названию и описанию..."
            value={q}
            onChange={(e) => setParam("q", e.target.value)}
          />
          {isAdmin && (
            <Checkbox
              label="Показать неактивных"
              checked={includeInactive}
              onChange={(e) => setParam("includeInactive", e.target.checked ? "true" : "")}
            />
          )}
        </div>
        <div className="assistants__categories">
          <Button
            variant={!categoryId ? "primary" : "secondary"}
            size="sm"
            onClick={() => setParam("categoryId", "")}
          >
            Все
          </Button>
          {categories.map((cat) => (
            <Button
              key={cat.id}
              variant={categoryId === cat.id ? "primary" : "secondary"}
              size="sm"
              onClick={() => setParam("categoryId", cat.id)}
            >
              {cat.name}
            </Button>
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