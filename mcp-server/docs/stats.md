# Статистика анкет

Эндпоинт `GET /stats` считает агрегаты по таблице `answers`.

- `language_counts` — `GROUP BY value` для `question_id = 2`
- `experience_counts` — `GROUP BY value` для `question_id = 3`
- `total_submissions` — `COUNT(*)` из `submissions`

Tool `get_stats` возвращает этот JSON как есть внутри поля `data`.

Если нужна одна анкета, используйте `get_submission` с конкретным `submission_id`.
Для справки по HTTP-контракту смотрите `docs/api.md`.
