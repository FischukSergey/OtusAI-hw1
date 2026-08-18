"""MCP-сервер мини-анкеты: транспорт stdio, регистрация tools через FastMCP."""

from typing import Any

from mcp.server.fastmcp import FastMCP

from tools_docs import lookup_docs_impl
from tools_survey import get_questions_impl, get_stats_impl, get_submission_impl

mcp = FastMCP("otusai-survey")


@mcp.tool()
def get_questions() -> dict[str, Any]:
    """Вернуть список вопросов мини-анкеты (GET /questions)."""
    return get_questions_impl()


@mcp.tool()
def get_stats() -> dict[str, Any]:
    """Вернуть агрегированную статистику по анкетам (GET /stats)."""
    return get_stats_impl()


@mcp.tool()
def get_submission(submission_id: int) -> dict[str, Any]:
    """Вернуть одну заполненную анкету по числовому id (GET /submissions/{id})."""
    return get_submission_impl(submission_id)


@mcp.tool()
def lookup_docs(topic: str) -> dict[str, Any]:
    """Найти справку по теме API/проекта в локальных markdown (только mcp-server/docs)."""
    return lookup_docs_impl(topic)


if __name__ == "__main__":
    mcp.run(transport="stdio")
