# Соответствие критериям оценки

Проверено e2e 2026-08-14 на LM Studio (`qwen/qwen3.6-35b-a3b`).  
Фактические логи: [`REPORT.md`](REPORT.md), [`e2e-logs/`](e2e-logs/).

## Агент запускается по инструкции

Инструкция: [`README.md`](README.md).

```bash
cd agent && source .venv/bin/activate && python main.py "покажи вопросы анкеты"
```

## Минимум один API-tool с реальным вызовом

**Объявление tool + HTTP-вызов:** [`tools_api.py`](tools_api.py)

| Что | Где |
|-----|-----|
| HTTP GET | `tools_api.py` L20–L25 (`_get`) |
| HTTP POST | `tools_api.py` L28–L33 (`_post`) |
| Tool `get_questions` | `tools_api.py` L36–L42 |
| Tool `submit_answers` | `tools_api.py` L45–L57 |
| Tool `get_stats` | `tools_api.py` L78–L84 |

**Debug print результата tool:** `tools_api.py` L15–L17 (`_debug`), видно в логах как `[tool:get_questions]`, `[tool:submit_answers]`, …

Пример из e2e:

```text
[tool:get_questions] [{"id": 1, "text": "Как вас зовут?", ...}]
[trace:AIMessage] tool_calls=[{'name': 'get_questions', ...}]
```

Файл: [`e2e-logs/01-get_questions.txt`](e2e-logs/01-get_questions.txt).

## Интерпретация запроса → выбор API-метода

| Запрос | Tool | Метод |
|--------|------|-------|
| покажи вопросы анкеты | `get_questions` | `GET /questions` |
| сохрани анкету (answers_json) | `submit_answers` | `POST /answers` |
| покажи статистику | `get_stats` | `GET /stats` |
| какой язык чаще? | `query_database` | NL→SQL |
| покажи список анкет | `list_submissions` | `GET /submissions` |

Подтверждение trace: [`REPORT.md`](REPORT.md), `[trace:AIMessage] tool_calls=...` в `e2e-logs/`.

## Контракт ответа

Описан в:

- [`prompts.py`](prompts.py) L31–L35
- [`PROMPTS.md`](PROMPTS.md)
- [`README.md`](README.md) → «Контракт ответа»

```text
Status: success | error
Action: ...
Data: ...
Errors: ...
```

## ≥5 проверочных запросов

[`REPORT.md`](REPORT.md) — запросы №1–6 с фактическими ответами.

## Промпты оформлены

[`PROMPTS.md`](PROMPTS.md) — system prompt, CLI-шаблоны, NL→SQL шаблон.

## Секреты не закоммичены

- `.env` в `.gitignore`
- в репозитории только [`.env.example`](.env.example) (`OPENAI_API_KEY=lm-studio`)
