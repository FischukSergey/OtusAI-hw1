"""Локальный прогон tools без Cursor: пишет [mcp]-логи в stderr.

Использование:
  SURVEY_API_BASE=http://127.0.0.1:8090 python smoke.py
"""

from __future__ import annotations

import json
import sys

from tools_docs import lookup_docs_impl
from tools_survey import get_questions_impl, get_stats_impl, get_submission_impl


def _dump(title: str, result: dict) -> None:
    print(f"=== {title} ===", file=sys.stderr)
    print(json.dumps(result, ensure_ascii=False, indent=2), file=sys.stderr)


def main() -> int:
    _dump("get_questions", get_questions_impl())
    _dump("get_stats", get_stats_impl())
    _dump("get_submission", get_submission_impl(1))
    _dump("lookup_docs", lookup_docs_impl("GET /stats"))
    _dump("lookup_docs empty", lookup_docs_impl("   "))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
