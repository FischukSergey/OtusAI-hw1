# MCP-сервер мини-анкеты

Собственный MCP-сервер для Cursor: четыре tool-а с JSON-результатом.
Три tool-а дублируют чтение REST API анкеты, четвёртый —
локальный doc-lookup в духе Context7.

## Принципы MCP

Cursor читает конфиг MCP (`.cursor/mcp.json`), запускает процесс
`command + args` и общается с ним по протоколу MCP через **stdio**.
Агент видит список tools: имя, описание, JSON Schema параметров — и сам
решает, какой tool вызвать из чата.

**Tool** в этом сервере — именованная функция с описанием, схемой входа и
структурированным JSON на выходе (`status` / `data` / `error`). Три tool-а
ходят только в фиксированные GET эндпоинты Go API, четвёртый читает только
файлы `mcp-server/docs/*.md`.

## Требования

1. Запущенный Go API: `http://127.0.0.1:8090` (см. корневой README).
2. Python 3.11+.
3. Cursor с включённым MCP.

Секреты не нужны. Единственная настройка — `SURVEY_API_BASE` в `.env.example`
и в `cursor-mcp.json`. SDK зафиксирован как `mcp==1.9.4` (FastMCP + stdio,
без тяжёлых зависимостей линейки 1.29+/2.x).

## Быстрый старт

```bash
cd backend
PORT=8090 go run .
```

```bash
cd mcp-server
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
```

Локальная проверка tools без Cursor (логи — в stderr):

```bash
SURVEY_API_BASE=http://127.0.0.1:8090 python smoke.py
```

## Как включить в Cursor

1. Поднимите Go API на `:8090`.
2. Создайте venv и поставьте зависимости (команды выше).
3. Скопируйте пример конфига в корень репозитория:
   `cp mcp-server/cursor-mcp.json .cursor/mcp.json`
   (папка `.cursor/` в gitignore — в репозитории лежит только пример).
4. Если путь к репозиторию другой — поправьте `command` и `args` в `.cursor/mcp.json`.
5. Cursor → Settings → MCP: сервер `otusai-survey` должен стать зелёным.
   В Agent-чате задайте запрос из [`REPORT.md`](REPORT.md).

Пример конфига: [`cursor-mcp.json`](cursor-mcp.json).

## Tools

| Tool | Параметры | Действие |
|------|-----------|----------|
| `get_questions` | нет | `GET /questions` |
| `get_stats` | нет | `GET /stats` |
| `get_submission` | `submission_id: int` | `GET /submissions/{id}` |
| `lookup_docs` | `topic: str` | поиск по `mcp-server/docs/*.md` |

## Контракт результата (Tool outputs contract)

Каждый tool возвращает объект, не голую строку:

```json
{
  "status": "success",
  "tool": "get_questions",
  "data": {},
  "error": null
}
```

Ошибка:

```json
{
  "status": "error",
  "tool": "get_submission",
  "data": null,
  "error": "HTTP 404: not found"
}
```

Для `lookup_docs` поле `data`:

```json
{
  "topic": "stats",
  "root": "mcp-server/docs",
  "matches": [
    {
      "file": "docs/stats.md",
      "title": "Статистика анкет",
      "excerpt": "...",
      "source": "mcp-server/docs/stats.md",
      "score": 3
    }
  ]
}
```

## Безопасность

- Нет shell, нет произвольного SQL, нет записи в БД.
- HTTP только на `SURVEY_API_BASE` и whitelist путей:
  `/questions`, `/stats`, `/submissions/{id}` (`tools_survey.py`).
- Файлы только из `mcp-server/docs/*.md` (`tools_docs.py`: `resolve` + `relative_to`).
- Логи без секретов: параметры — `topic` / `submission_id`.
- Логи пишутся в **stderr**: stdout занят MCP stdio.

## Подтверждения ссылками на код

### Реализован MCP-сервер

Подъём сервера и регистрация tools: [`server.py`](server.py) L10–L38.

- `FastMCP("otusai-survey")` — L10
- декораторы `@mcp.tool()` — L13–L34
- `mcp.run(transport="stdio")` — L37–L38

### Реализованы инструменты

| Tool | Объявление | Реализация | Debug-лог |
|------|------------|------------|-----------|
| `get_questions` | `server.py` L13–L16 | `tools_survey.py` L44–L50 | `logutil.py` L11–L26, вызов `tools_survey.py` L20–L27 |
| `get_stats` | `server.py` L19–L22 | `tools_survey.py` L53–L59 | то же |
| `get_submission` | `server.py` L25–L28 | `tools_survey.py` L62–L72 | то же |
| `lookup_docs` | `server.py` L31–L34 | `tools_docs.py` L69–L103 | `logutil.py` L11–L26, вызов `tools_docs.py` L16–L23 |

Пример вывода (stderr):

```text
[mcp] {"ts": "2026-08-18T18:49:43.552456+00:00", "tool": "get_questions", "params": {}, "status": "success"}
[mcp] {"ts": "2026-08-18T18:49:43.439987+00:00", "tool": "get_submission", "params": {"submission_id": 1}, "status": "success"}
[mcp] {"ts": "2026-08-18T18:49:43.510831+00:00", "tool": "lookup_docs", "params": {"topic": "GET /stats"}, "status": "success"}
```

Источник: Cursor Agent, 2026-08-18 21:49. Сводка: [`REPORT.md`](REPORT.md).

### Агент вызывает tool

Запрос: `Покажи вопросы анкеты`  
Ожидаемый tool: `get_questions`  
Фактическое подтверждение:

- stderr сервера + `CallToolRequest` в Cursor:
  [`logs/01-get_questions.txt`](logs/01-get_questions.txt)
- полный trace четырёх вызовов:
  [`logs/cursor-agent-stdio.log`](logs/cursor-agent-stdio.log)
- JSON-ответ tool:
  [`logs/cursor-results/01-get_questions.json`](logs/cursor-results/01-get_questions.json)

Тот же прогон: `get_stats`, `get_submission`, `lookup_docs` — см. [`REPORT.md`](REPORT.md).

### Контракт результата

Этот файл, раздел **«Контракт результата (Tool outputs contract)»** выше.

## Секреты

- Реальные `.env` не коммитятся (корневой `.gitignore`).
- В репозитории — [`.env.example`](.env.example) с `SURVEY_API_BASE`.
