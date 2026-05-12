import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { getMyRuns } from "@/api/runs";
import RunRow from "@/components/RunRow";
import Pagination from "@/components/Pagination";
import Spinner from "@/components/ui/Spinner";
import Empty from "@/components/ui/Empty";
import "./MyRuns.css";

const STATUSES = [
  { value: "", label: "Все" },
  { value: "pending", label: "В обработке" },
  { value: "success", label: "Успешные" },
  { value: "failed", label: "Ошибки" },
];

export default function MyRuns() {
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get("page") || "1");
  const status = searchParams.get("status") || "";

  const [runs, setRuns] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, pageSize: 20, total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    setLoading(true);
    setError(null);

    const params = { page, pageSize: 20 };
    if (status) params.status = status;

    getMyRuns(params)
      .then(({ runs, pagination }) => {
        setRuns(runs);
        setPagination(pagination);
      })
      .catch(() => setError("Не удалось загрузить историю запусков"))
      .finally(() => setLoading(false));
  }, [page, status]);

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
    <div className="my-runs">
      <div className="my-runs__header">
        <h1 className="my-runs__title">Мои запуски</h1>
      </div>

      <div className="my-runs__filters">
        {STATUSES.map((s) => (
          <button
            key={s.value}
            className={`my-runs__filter-btn ${status === s.value ? "my-runs__filter-btn--active" : ""}`}
            onClick={() => setParam("status", s.value)}
          >
            {s.label}
          </button>
        ))}
      </div>

      {loading && (
        <div className="my-runs__spinner">
          <Spinner size="lg" />
        </div>
      )}

      {error && !loading && (
        <p className="my-runs__error">{error}</p>
      )}

      {!loading && !error && runs.length === 0 && (
        <Empty
          title="Запусков пока нет"
          description="Откройте любого ассистента и запустите его"
        />
      )}

      {!loading && !error && runs.length > 0 && (
        <>
          <div className="my-runs__table-wrap">
            <table className="my-runs__table">
              <thead>
                <tr>
                  <th>Ассистент</th>
                  <th>Запрос</th>
                  <th>Статус</th>
                  <th>Дата</th>
                  <th></th>
                </tr>
              </thead>
              <tbody>
                {runs.map((run) => (
                  <RunRow key={run.id} run={run} />
                ))}
              </tbody>
            </table>
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