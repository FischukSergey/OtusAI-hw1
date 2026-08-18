# REST API мини-анкеты

Базовый URL локальной разработки: `http://127.0.0.1:8090`.

Сервер — Go + SQLite. Запись ответов идёт только через `POST /answers`.
MCP-сервер вызывает **только безопасные GET**: `/questions`, `/stats`, `/submissions/{id}`.

## GET /questions

Список вопросов анкеты.

Каждый элемент:

- `id` — числовой идентификатор вопроса
- `text` — формулировка
- `type` — `text` | `radio` | `select`
- `options` — варианты ответа (для radio/select)

## POST /answers

Сохранить заполненную анкету.

Тело:

```json
{
  "answers": [
    {"question_id": 1, "value": "Иван"},
    {"question_id": 2, "value": "Go"}
  ]
}
```

Ответ: `{"status":"ok","id":1}`.

MCP-сервер этот метод **не вызывает** (нет write-tool).

## GET /submissions

Список всех заполненных анкет: `id`, `created_at`, `answers`.

## GET /submissions/{id}

Одна анкета по id. 404, если записи нет.

## GET /stats

Агрегаты:

- `total_submissions` — число анкет
- `language_counts` — частоты ответа на вопрос 2 (язык)
- `experience_counts` — частоты ответа на вопрос 3 (опыт)
