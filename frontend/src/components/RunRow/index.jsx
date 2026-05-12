import './RunRow.css';
import { useState } from "react";
import StatusBadge from "@/components/StatusBadge";
import Button from "@/components/ui/Button";

export default function RunRow({ run }) {
  const [expanded, setExpanded] = useState(false);

  return (
    <>
      <tr className="run-row">
        <td className="run-row__cell">
          <span className="run-row__assistant">{run.assistantName}</span>
          <span className="run-row__category">{run.categoryName}</span>
        </td>
        <td className="run-row__cell run-row__prompt">
          {run.userPrompt}
        </td>
        <td className="run-row__cell">
          <StatusBadge status={run.status} />
        </td>
        <td className="run-row__cell run-row__date">
          {new Date(run.createdAt).toLocaleString("ru-RU")}
        </td>
        <td className="run-row__cell">
          {(run.output || run.error) && (
            <Button
              variant="ghost"
              size="sm"
              onClick={() => setExpanded((v) => !v)}
            >
              {expanded ? "Скрыть" : "Показать"}
            </Button>
          )}
        </td>
      </tr>
      {expanded && (
        <tr className="run-row__expanded">
          <td colSpan={5}>
            {run.output && (
              <pre className="run-row__output">{run.output}</pre>
            )}
            {run.error && (
              <p className="run-row__error">{run.error}</p>
            )}
          </td>
        </tr>
      )}
    </>
  );
}