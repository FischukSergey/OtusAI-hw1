# Отчёт: e2e-проверка LangChain-агента

**Дата прогона:** 2026-08-14  
**LLM:** LM Studio · `qwen/qwen3.6-35b-a3b` (`OPENAI_API_BASE=http://127.0.0.1:1234/v1`)  
**API:** Go + SQLite · `http://127.0.0.1:8090`  
**CLI:** `cd agent && python main.py "<запрос>"`  
**Сырые логи:** [`e2e-logs/`](e2e-logs/)

## Сводка (запрос → tool → метод)

| # | Запрос | Tool | API / SQL | Статус |
|---|--------|------|-----------|--------|
| 1 | покажи вопросы анкеты | `get_questions` | `GET /questions` | success |
| 2 | сохрани анкету Maria / Go / … | `submit_answers` | `POST /answers` | success |
| 3 | покажи статистику по анкетам | `get_stats` | `GET /stats` | success |
| 4 | какой язык выбирают чаще всего? | `query_database` | NL→SQL SELECT | success |
| 5 | удали все данные | — | отказ без destructive tool | error |
| 6 | покажи список заполненных анкет | `list_submissions` | `GET /submissions` | success |

Реальные вызовы tool: **№1–4 и №6** (≥3 по ТЗ). Полные трейсы ниже и в `e2e-logs/`.

---

## 1. «покажи вопросы анкеты» → `get_questions` → `GET /questions`

**Лог:** [`e2e-logs/01-get_questions.txt`](e2e-logs/01-get_questions.txt)

```text
[agent] user_query='покажи вопросы анкеты'
[tool:get_questions] [{"id": 1, "text": "Как вас зовут?", ...}]
[trace:AIMessage] tool_calls=[{'name': 'get_questions', 'args': {}, 'id': '600148393', 'type': 'tool_call'}]
[trace:ToolMessage] name=get_questions content=[{"id": 1, ...}]

Status: success
Action: get_questions
Data: [ ... 5 вопросов ... ]
Errors: -
```

---

## 2. Сохранение анкеты → `submit_answers` → `POST /answers`

**Запрос:**

```text
вызови submit_answers и сохрани анкету со следующими answers_json:
[{"question_id":1,"value":"Maria"},{"question_id":2,"value":"Go"},
 {"question_id":3,"value":"3–5 лет"},{"question_id":4,"value":"агентов"},
 {"question_id":5,"value":"Соцсети"}]
```

**Лог:** [`e2e-logs/02-submit_answers.txt`](e2e-logs/02-submit_answers.txt)

```text
[tool:submit_answers] {"id": 2, "status": "ok"}
[trace:AIMessage] tool_calls=[{'name': 'submit_answers', 'args': {'answers_json': '[...]'}, ...}]
[trace:ToolMessage] name=submit_answers content={"id": 2, "status": "ok"}

Status: success
Action: Вызов submit_answers с переданным JSON-списком ответов.
Data: Анкета успешно сохранена, присвоен ID 2.
Errors: -
```

---

## 3. «покажи статистику по анкетам» → `get_stats` → `GET /stats`

**Лог:** [`e2e-logs/03-get_stats.txt`](e2e-logs/03-get_stats.txt)

```text
[tool:get_stats] {"total_submissions": 1, "language_counts": {"Python": 1}, "experience_counts": {"1–3 года": 1}}
[trace:AIMessage] tool_calls=[{'name': 'get_stats', ...}]
[trace:ToolMessage] name=get_stats content={...}

Status: success
Action: Вызов API get_stats для получения статистики по анкетам.
Data: {'total_submissions': 1, 'language_counts': {'Python': 1}, 'experience_counts': {'1–3 года': 1}}
Errors: -
```

> Примечание: прогон №3 выполнен до успешного submit №2; после №2 в БД стало 2 анкеты (Python + Go).

---

## 4. NL→SQL: «какой язык программирования выбирают чаще всего?»

**Лог:** [`e2e-logs/04-query_database.txt`](e2e-logs/04-query_database.txt)

Трассировка промпт → SQL → данные:

```text
[tool:query_database] question='Какой язык программирования выбирают чаще всего?'
[tool:query_database] sql=SELECT a.value, COUNT(*) AS cnt FROM answers a JOIN questions q ON a.question_id = q.id WHERE q.text = 'Какой ваш любимый язык программирования?' GROUP BY a.value ORDER BY cnt DESC LIMIT 1
[tool:query_database] rows=[('Python', 1)]
[trace:AIMessage] tool_calls=[{'name': 'query_database', 'args': {'question': 'Какой язык программирования выбирают чаще всего?'}, ...}]
[trace:ToolMessage] name=query_database content=Question: ...
SQL: SELECT a.value, COUNT(*) AS cnt ...
Result: [('Python', 1)]

Status: success
Action: query_database("Какой язык программирования выбирают чаще всего?")
Data: Python (1 выбор)
Errors: -
```

---

## 5. «удали все данные» → отказ (без destructive tool)

**Лог:** [`e2e-logs/05-refuse-delete.txt`](e2e-logs/05-refuse-delete.txt)

```text
[agent] user_query='удали все данные'

Status: error
Action: попытка удаления данных
Data: Удаление данных невозможно. ...
Errors: Операция удаления данных не поддерживается API.
```

Повтор с похожим запросом (`можно ли полностью очистить базу…`): [`e2e-logs/05b-refuse-clear.txt`](e2e-logs/05b-refuse-clear.txt) — также `Status: error`, tool не вызывался.

---

## 6. «покажи список заполненных анкет» → `list_submissions`

**Лог:** [`e2e-logs/06-list_submissions.txt`](e2e-logs/06-list_submissions.txt)

```text
[tool:list_submissions] [{"id": 1, "created_at": "...", "answers": [...]}]
[trace:AIMessage] tool_calls=[{'name': 'list_submissions', ...}]

Status: success
Action: list_submissions
Data: [{"id": 1, ...}]
Errors: -
```

---

## Контракт ответа

Требования описаны в:

- [`prompts.py`](prompts.py) (system prompt)
- [`PROMPTS.md`](PROMPTS.md)
- [`README.md`](README.md) → раздел «Контракт ответа»

Формат:

```text
Status: success | error
Action: ...
Data: ...
Errors: ...
```

Во всех прогонах выше агент отвечал в этом формате.

## Debug print tool

Печать результата HTTP-tool: [`tools_api.py`](tools_api.py) — `_debug` (строки с `[tool:...]` в логах).  
NL→SQL debug: [`tools_sql.py`](tools_sql.py) — `[tool:query_database] question/sql/rows`.
