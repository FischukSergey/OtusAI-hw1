#!/usr/bin/env python3
"""CLI: python main.py \"покажи вопросы анкеты\""""

from __future__ import annotations

import sys

from dotenv import load_dotenv

load_dotenv()

from agent import run_agent  # noqa: E402


def main() -> int:
    if len(sys.argv) < 2 or not sys.argv[1].strip():
        print('Использование: python main.py "<запрос на естественном языке>"', file=sys.stderr)
        print('Пример: python main.py "покажи вопросы анкеты"', file=sys.stderr)
        return 2

    query = " ".join(sys.argv[1:]).strip()
    try:
        answer = run_agent(query)
    except Exception as exc:  # noqa: BLE001 — CLI показывает ошибку пользователю
        print(
            "Status: error\n"
            f"Action: agent failed\n"
            f"Data: -\n"
            f"Errors: {exc}",
            file=sys.stderr,
        )
        return 1

    print(answer)
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
