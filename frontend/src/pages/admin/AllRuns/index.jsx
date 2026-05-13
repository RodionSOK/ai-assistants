import { useState, useEffect } from "react";
import { useSearchParams } from "react-router-dom";
import { getAllRuns } from "@/api/runs";
import { getAssistants } from "@/api/assistants";
import RunRow from "@/components/RunRow";
import Pagination from "@/components/Pagination";
import Spinner from "@/components/ui/Spinner";
import Empty from "@/components/ui/Empty";
import Button from "@/components/ui/Button";
import Select from "@/components/ui/Select";
import "./AllRuns.css";

const STATUSES = [
  { value: "", label: "Все" },
  { value: "pending", label: "В обработке" },
  { value: "success", label: "Успешные" },
  { value: "failed", label: "Ошибки" },
];

export default function AllRuns() {
  const [searchParams, setSearchParams] = useSearchParams();

  const page = parseInt(searchParams.get("page") || "1");
  const status = searchParams.get("status") || "";
  const assistantId = searchParams.get("assistantId") || "";

  const [runs, setRuns] = useState([]);
  const [pagination, setPagination] = useState({ page: 1, pageSize: 20, total: 0 });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [assistants, setAssistants] = useState([]);

  useEffect(() => {
    getAssistants({ pageSize: 100, includeInactive: true })
      .then(({ assistants }) => setAssistants(assistants))
      .catch(() => {});
  }, []);

  useEffect(() => {
    setLoading(true);
    setError(null);

    const params = { page, pageSize: 20 };
    if (status) params.status = status;
    if (assistantId) params.assistantId = assistantId;

    getAllRuns(params)
      .then(({ runs, pagination }) => {
        setRuns(runs);
        setPagination(pagination);
      })
      .catch(() => setError("Не удалось загрузить запуски"))
      .finally(() => setLoading(false));
  }, [page, status, assistantId]);

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
    <div className="all-runs">
      <div className="all-runs__header">
        <h1 className="all-runs__title">Все запуски</h1>
      </div>

      <div className="all-runs__filters">
        <Select
          placeholder="Все ассистенты"
          value={assistantId}
          onChange={(e) => setParam("assistantId", e.target.value)}
          options={assistants.map((a) => ({ value: a.id, label: a.name }))}
        />
        {(assistantId || status) && (
          <Button
            variant="ghost"
            size="sm"
            onClick={() => {
              setSearchParams((prev) => {
                const next = new URLSearchParams(prev);
                next.delete("assistantId");
                next.delete("status");
                next.delete("page");
                return next;
              });
            }}
          >
            Сбросить фильтры
          </Button>
        )}

        <div className="all-runs__status-filters">
          {STATUSES.map((s) => (
            <Button
              key={s.value}
              variant={status === s.value ? "primary" : "secondary"}
              size="sm"
              onClick={() => setParam("status", s.value)}
            >
              {s.label}
            </Button>
          ))}
        </div>
      </div>

      {loading && (
        <div className="all-runs__spinner">
          <Spinner size="lg" />
        </div>
      )}

      {error && !loading && (
        <p className="all-runs__error">{error}</p>
      )}

      {!loading && !error && runs.length === 0 && (
        <Empty
          title="Запусков нет"
          description="Пока никто не запускал ассистентов"
        />
      )}

      {!loading && !error && runs.length > 0 && (
        <>
          <div className="all-runs__table-wrap">
            <table className="all-runs__table">
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