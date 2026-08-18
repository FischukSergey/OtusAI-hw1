"""Локальный doc-lookup (аналог Context7): только markdown в mcp-server/docs/."""

from __future__ import annotations

import re
from pathlib import Path
from typing import Any

from logutil import log_tool

DOCS_ROOT = (Path(__file__).resolve().parent / "docs").resolve()
_MAX_EXCERPT = 400
_MAX_MATCHES = 5


def _ok(tool: str, params: dict[str, Any], data: Any) -> dict[str, Any]:
    log_tool(tool, params, "success")
    return {"status": "success", "tool": tool, "data": data, "error": None}


def _err(tool: str, params: dict[str, Any], error: str) -> dict[str, Any]:
    log_tool(tool, params, "error", error)
    return {"status": "error", "tool": tool, "data": None, "error": error}


def _safe_docs() -> list[Path]:
    if not DOCS_ROOT.is_dir():
        return []
    files: list[Path] = []
    for path in sorted(DOCS_ROOT.glob("*.md")):
        resolved = path.resolve()
        resolved.relative_to(DOCS_ROOT)
        if resolved.is_file():
            files.append(resolved)
    return files


def _tokens(topic: str) -> list[str]:
    return [t for t in re.split(r"[^\w]+", topic.lower(), flags=re.UNICODE) if t]


def _excerpt(text: str, topic: str) -> str:
    lowered = text.lower()
    idx = lowered.find(topic.lower())
    if idx < 0:
        for token in _tokens(topic):
            idx = lowered.find(token)
            if idx >= 0:
                break
    if idx < 0:
        idx = 0
    start = max(0, idx - 80)
    end = min(len(text), idx + _MAX_EXCERPT)
    snippet = text[start:end].strip()
    if start > 0:
        snippet = "…" + snippet
    if end < len(text):
        snippet = snippet + "…"
    return snippet


def _title(text: str, filename: str) -> str:
    for line in text.splitlines():
        if line.startswith("#"):
            return line.lstrip("#").strip()
    return filename


def lookup_docs_impl(topic: str) -> dict[str, Any]:
    params: dict[str, Any] = {"topic": topic}
    try:
        query = (topic or "").strip()
        if not query:
            return _err("lookup_docs", params, "topic не должен быть пустым")

        tokens = _tokens(query)
        matches: list[dict[str, Any]] = []
        for path in _safe_docs():
            text = path.read_text(encoding="utf-8")
            hay = f"{path.name}\n{text}".lower()
            score = sum(hay.count(token) for token in tokens) if tokens else 0
            if score == 0 and query.lower() not in hay:
                continue
            rel = f"docs/{path.name}"
            matches.append(
                {
                    "file": rel,
                    "title": _title(text, path.name),
                    "excerpt": _excerpt(text, query),
                    "source": f"mcp-server/{rel}",
                    "score": score,
                }
            )

        matches.sort(key=lambda m: m["score"], reverse=True)
        data = {
            "topic": query,
            "root": "mcp-server/docs",
            "matches": matches[:_MAX_MATCHES],
        }
        return _ok("lookup_docs", params, data)
    except Exception as exc:
        return _err("lookup_docs", params, str(exc))
