"""LangChain NL→SQL tool (read-only) against the same SQLite DB as the Go API."""

from __future__ import annotations

import os
import re
from pathlib import Path

from langchain_community.utilities import SQLDatabase
from langchain_core.language_models.chat_models import BaseChatModel
from langchain_core.tools import StructuredTool
from pydantic import BaseModel, Field

FORBIDDEN_SQL = re.compile(
    r"\b(INSERT|UPDATE|DELETE|DROP|ALTER|CREATE|REPLACE|ATTACH|DETACH|PRAGMA|VACUUM)\b",
    re.IGNORECASE,
)


class NLQueryInput(BaseModel):
    question: str = Field(description="Вопрос на естественном языке к данным анкеты в SQLite")


def _resolve_db_path() -> str:
    raw = os.getenv("DB_PATH", "../backend/data/survey.db")
    path = Path(raw)
    if not path.is_absolute():
        path = (Path(__file__).resolve().parent / path).resolve()
    return str(path)


def _assert_select_only(sql: str) -> str:
    cleaned = sql.strip().rstrip(";")
    # model sometimes wraps SQL in ```sql ... ```
    cleaned = re.sub(r"^```(?:sql)?\s*", "", cleaned, flags=re.IGNORECASE)
    cleaned = re.sub(r"\s*```$", "", cleaned)
    cleaned = cleaned.strip().rstrip(";")

    if FORBIDDEN_SQL.search(cleaned):
        raise ValueError(f"Только SELECT разрешён, получен: {cleaned}")
    if not cleaned.lower().lstrip().startswith("select"):
        raise ValueError(f"Ожидался SELECT, получен: {cleaned}")
    return cleaned


def build_nl_sql_tool(llm: BaseChatModel) -> StructuredTool:
    db_path = _resolve_db_path()
    db = SQLDatabase.from_uri(
        f"sqlite:///{db_path}",
        include_tables=["questions", "submissions", "answers"],
        sample_rows_in_table_info=2,
    )
    schema = db.get_table_info()

    def run_nl_query(question: str) -> str:
        prompt = (
            "Ты генерируешь SQL для SQLite.\n"
            "Правила:\n"
            "- Верни ТОЛЬКО один SQL-запрос SELECT, без пояснений.\n"
            "- Нельзя INSERT/UPDATE/DELETE/DROP/ALTER/CREATE/PRAGMA.\n"
            "- Используй только таблицы: questions, submissions, answers.\n\n"
            f"Схема:\n{schema}\n\n"
            f"Вопрос: {question}\n"
            "SQL:"
        )
        raw = llm.invoke(prompt).content
        if isinstance(raw, list):
            raw = "".join(
                part.get("text", "") if isinstance(part, dict) else str(part) for part in raw
            )
        sql = _assert_select_only(str(raw))
        print(f"[tool:query_database] question={question!r}", flush=True)
        print(f"[tool:query_database] sql={sql}", flush=True)
        rows = db.run(sql)
        print(f"[tool:query_database] rows={rows}", flush=True)
        return (
            f"Question: {question}\n"
            f"SQL: {sql}\n"
            f"Result: {rows}"
        )

    return StructuredTool.from_function(
        func=run_nl_query,
        name="query_database",
        description=(
            "Аналитический вопрос к SQLite (NL→SQL). "
            "Используй для агрегатов и выборок, когда нужен произвольный SELECT. "
            "Для записи анкеты используй submit_answers, не этот tool."
        ),
        args_schema=NLQueryInput,
    )
