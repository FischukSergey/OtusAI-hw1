# Отчёт: проверка MCP-сервера

**Дата прогона:** 2026-08-18 21:49  
**Клиент:** Cursor Agent · этот чат · сервер `project-0-OtusAI-hw1-otusai-survey`  
**Сервер:** `mcp-server/server.py` · FastMCP `otusai-survey` · stdio  
**API:** Go + SQLite · `http://127.0.0.1:8090`  
**Trace Cursor:** [`logs/cursor-agent-stdio.log`](logs/cursor-agent-stdio.log)  
**JSON-ответы tools:** [`logs/cursor-results/`](logs/cursor-results/)

Это не `smoke.py`: агент IDE сам вызвал MCP-tools. Cursor записал
`CallToolRequest` + `[mcp]` в свой MCP-лог; копия лежит в репозитории.

## Сводка (запрос → tool)

| # | Запрос | Ожидаемый tool | Вызов MCP | Статус |
|---|--------|----------------|-----------|--------|
| 1 | Покажи вопросы анкеты | `get_questions` | да | success |
| 2 | Какая статистика по анкетам? | `get_stats` | да | success |
| 3 | Покажи анкету с id 1 | `get_submission` | да | success |
| 4 | Найди в документации, как устроен GET /stats | `lookup_docs` | да | success |
| 5 | Удали все данные анкеты | — | нет | отказ без destructive tool |

Реальные вызовы из IDE: **№1–4**.

---

## 1. «Покажи вопросы анкеты» → `get_questions`

**Ожидаемый tool:** `get_questions`  
**Подтверждение:** [`logs/01-get_questions.txt`](logs/01-get_questions.txt),
полный ответ [`logs/cursor-results/01-get_questions.json`](logs/cursor-results/01-get_questions.json)

```text
2026-08-18 21:49:43.546 Processing request of type CallToolRequest
2026-08-18 21:49:43.552 HTTP Request: GET http://127.0.0.1:8090/questions "HTTP/1.1 200 OK"
[mcp] {"ts": "2026-08-18T18:49:43.552456+00:00", "tool": "get_questions", "params": {}, "status": "success"}
```

`data`: 5 вопросов (`id` 1–5).

---

## 2. «Какая статистика по анкетам?» → `get_stats`

**Ожидаемый tool:** `get_stats`  
**Подтверждение:** [`logs/02-get_stats.txt`](logs/02-get_stats.txt),
[`logs/cursor-results/02-get_stats.json`](logs/cursor-results/02-get_stats.json)

```text
2026-08-18 21:49:45.045 Processing request of type CallToolRequest
2026-08-18 21:49:45.056 HTTP Request: GET http://127.0.0.1:8090/stats "HTTP/1.1 200 OK"
[mcp] {"ts": "2026-08-18T18:49:45.056512+00:00", "tool": "get_stats", "params": {}, "status": "success"}
```

`data`: `total_submissions=2`, языки Go/Python.

---

## 3. «Покажи анкету с id 1» → `get_submission`

**Ожидаемый tool:** `get_submission`  
**Подтверждение:** [`logs/03-get_submission.txt`](logs/03-get_submission.txt),
[`logs/cursor-results/03-get_submission.json`](logs/cursor-results/03-get_submission.json)

```text
2026-08-18 21:49:43.410 Processing request of type CallToolRequest
2026-08-18 21:49:43.440 HTTP Request: GET http://127.0.0.1:8090/submissions/1 "HTTP/1.1 200 OK"
[mcp] {"ts": "2026-08-18T18:49:43.439987+00:00", "tool": "get_submission", "params": {"submission_id": 1}, "status": "success"}
```

`data.id=1`, ответы Alex / Python / LangChain.

---

## 4. «Найди в документации, как устроен GET /stats» → `lookup_docs`

**Ожидаемый tool:** `lookup_docs`  
**Подтверждение:** [`logs/04-lookup_docs.txt`](logs/04-lookup_docs.txt),
[`logs/cursor-results/04-lookup_docs.json`](logs/cursor-results/04-lookup_docs.json)

```text
2026-08-18 21:49:43.510 Processing request of type CallToolRequest
[mcp] {"ts": "2026-08-18T18:49:43.510831+00:00", "tool": "lookup_docs", "params": {"topic": "GET /stats"}, "status": "success"}
```

`data.matches`: `docs/api.md` (score 7), `docs/stats.md` (score 6).

---

## 5. «Удали все данные анкеты» → без tool

У сервера нет write/delete tools. Агент не вызывает MCP-tool.

Файл: [`logs/05-refuse-delete.txt`](logs/05-refuse-delete.txt).
