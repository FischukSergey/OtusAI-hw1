# Соответствие критериям ДЗ (MCP-сервер)

## Сервер запускается по инструкции

Инструкция: [`README.md`](README.md).

```bash
cd mcp-server && source .venv/bin/activate
SURVEY_API_BASE=http://127.0.0.1:8090 python smoke.py
```

В Cursor: скопировать [`cursor-mcp.json`](cursor-mcp.json) → `.cursor/mcp.json`.

## Минимум 2 инструмента (здесь 4)

| Tool | Файл / строки |
|------|----------------|
| `get_questions` | `server.py` L13–L16, `tools_survey.py` L44–L50 |
| `get_stats` | `server.py` L19–L22, `tools_survey.py` L53–L59 |
| `get_submission` | `server.py` L25–L28, `tools_survey.py` L62–L72 |
| `lookup_docs` | `server.py` L31–L34, `tools_docs.py` L69–L103 |

Описания и схемы — docstring + type hints FastMCP (`submission_id: int`, `topic: str`).

## Интеграция с агентом в IDE

Пример конфига: [`cursor-mcp.json`](cursor-mcp.json).  
Шаги «как включить»: [`README.md`](README.md) → «Как включить в Cursor».

## ≥5 проверочных запросов, ≥3 с вызовом MCP-tool

[`REPORT.md`](REPORT.md): запросы №1–4 вызваны Cursor Agent 2026-08-18 21:49.
Trace: [`logs/cursor-agent-stdio.log`](logs/cursor-agent-stdio.log).
№5 — отказ без destructive tool.

## Логи / отладка на стороне сервера

[`logutil.py`](logutil.py) L11–L26 — имя tool, params, `status`.  
Вызовы: `tools_survey.py` L20–L27, `tools_docs.py` L16–L23.  
Примеры: [`logs/`](logs/).

## Секреты не закоммичены

- `.env` в корневом `.gitignore`
- в репозитории только [`.env.example`](.env.example)

## Контракт результата

[`README.md`](README.md) → «Контракт результата (Tool outputs contract)».
