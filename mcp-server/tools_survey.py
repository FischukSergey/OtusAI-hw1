"""HTTP-tools к REST API мини-анкеты. Только фиксированные GET-пути."""

from __future__ import annotations

import os
from typing import Any

import httpx

from logutil import log_tool

API_BASE = os.getenv("SURVEY_API_BASE", "http://127.0.0.1:8090").rstrip("/")

_ALLOWED_GET = (
    "/questions",
    "/stats",
)


def _ok(tool: str, params: dict[str, Any], data: Any) -> dict[str, Any]:
    log_tool(tool, params, "success")
    return {"status": "success", "tool": tool, "data": data, "error": None}


def _err(tool: str, params: dict[str, Any], error: str) -> dict[str, Any]:
    log_tool(tool, params, "error", error)
    return {"status": "error", "tool": tool, "data": None, "error": error}


def _get(path: str) -> Any:
    if path not in _ALLOWED_GET and not path.startswith("/submissions/"):
        raise ValueError(f"путь не в whitelist: {path}")
    if path.startswith("/submissions/"):
        rest = path.removeprefix("/submissions/")
        if not rest.isdigit() or int(rest) <= 0:
            raise ValueError(f"некорректный submission_id в пути: {path}")
    url = f"{API_BASE}{path}"
    with httpx.Client(timeout=30.0) as client:
        resp = client.get(url)
        resp.raise_for_status()
        return resp.json()


def get_questions_impl() -> dict[str, Any]:
    params: dict[str, Any] = {}
    try:
        data = _get("/questions")
        return _ok("get_questions", params, data)
    except Exception as exc:
        return _err("get_questions", params, str(exc))


def get_stats_impl() -> dict[str, Any]:
    params: dict[str, Any] = {}
    try:
        data = _get("/stats")
        return _ok("get_stats", params, data)
    except Exception as exc:
        return _err("get_stats", params, str(exc))


def get_submission_impl(submission_id: int) -> dict[str, Any]:
    params: dict[str, Any] = {"submission_id": submission_id}
    try:
        if submission_id <= 0:
            return _err("get_submission", params, "submission_id должен быть > 0")
        data = _get(f"/submissions/{submission_id}")
        return _ok("get_submission", params, data)
    except httpx.HTTPStatusError as exc:
        return _err("get_submission", params, f"HTTP {exc.response.status_code}: {exc.response.text}")
    except Exception as exc:
        return _err("get_submission", params, str(exc))
