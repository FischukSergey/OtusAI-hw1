"""LangChain tools that call the Survey REST API over HTTP."""

from __future__ import annotations

import json
import os
from typing import Any

import httpx
from langchain_core.tools import tool

API_BASE = os.getenv("SURVEY_API_BASE", "http://127.0.0.1:8090").rstrip("/")


def _debug(tool_name: str, payload: Any) -> None:
    """Печать результата tool в консоль для дебага (критерий ДЗ)."""
    print(f"[tool:{tool_name}] {payload}", flush=True)


def _get(path: str) -> Any:
    url = f"{API_BASE}{path}"
    with httpx.Client(timeout=30.0) as client:
        resp = client.get(url)
        resp.raise_for_status()
        return resp.json()


def _post(path: str, body: dict[str, Any]) -> Any:
    url = f"{API_BASE}{path}"
    with httpx.Client(timeout=30.0) as client:
        resp = client.post(url, json=body)
        resp.raise_for_status()
        return resp.json()


@tool
def get_questions() -> str:
    """Получить список вопросов анкеты (GET /questions)."""
    data = _get("/questions")
    result = json.dumps(data, ensure_ascii=False)
    _debug("get_questions", result)
    return result


@tool
def submit_answers(answers_json: str) -> str:
    """Сохранить ответы анкеты (POST /answers).

    Аргумент answers_json — JSON-строка вида:
    [{"question_id": 1, "value": "Alex"}, {"question_id": 2, "value": "Python"}, ...]
    Нужны ответы на все вопросы анкеты.
    """
    answers = json.loads(answers_json)
    data = _post("/answers", {"answers": answers})
    result = json.dumps(data, ensure_ascii=False)
    _debug("submit_answers", result)
    return result


@tool
def list_submissions() -> str:
    """Получить список всех заполненных анкет (GET /submissions)."""
    data = _get("/submissions")
    result = json.dumps(data, ensure_ascii=False)
    _debug("list_submissions", result)
    return result


@tool
def get_submission(submission_id: int) -> str:
    """Получить одну анкету по id (GET /submissions/{id})."""
    data = _get(f"/submissions/{submission_id}")
    result = json.dumps(data, ensure_ascii=False)
    _debug("get_submission", result)
    return result


@tool
def get_stats() -> str:
    """Получить статистику по анкетам (GET /stats): число submissions и агрегаты."""
    data = _get("/stats")
    result = json.dumps(data, ensure_ascii=False)
    _debug("get_stats", result)
    return result


API_TOOLS = [
    get_questions,
    submit_answers,
    list_submissions,
    get_submission,
    get_stats,
]
