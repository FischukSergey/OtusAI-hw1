# LangChain-агент: NL-обёртка над API анкеты

Минимальный AI-агент на LangChain/LangGraph, который принимает запрос на естественном языке,
вызывает REST API мини-анкеты через tools и (для аналитики) превращает вопрос в SQL к SQLite.

## Требования

1. Запущенный Go API (`http://127.0.0.1:8090`) с SQLite.
2. [LM Studio](https://lmstudio.ai/) с загруженной моделью и Local Server на порту `1234`
   (OpenAI-compatible API).
3. Python 3.11+.

> Модели из чата Cursor **нельзя** вызвать из `python main.py`. Агенту нужен свой endpoint —
> здесь это LM Studio.

## Быстрый старт

### 1. Поднять API

```bash
cd backend
PORT=8090 go run .
# SQLite: ./data/survey.db
# http://localhost:8090
```

Или Docker (снаружи порт **8090**):

```bash
docker compose up --build
# http://localhost:8090
```

### 2. Запустить LM Studio

1. Загрузить любую instruction-модель (желательно с tool/function calling).
2. Developer → Start Server (`http://127.0.0.1:1234`).
3. Скопировать id модели в `OPENAI_MODEL`.

### 3. Настроить агента

```bash
cd agent
python3 -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
cp .env.example .env
# при необходимости поправьте OPENAI_MODEL и пути
```

### 4. Запуск CLI

```bash
python main.py "покажи вопросы анкеты"
python main.py "какой язык программирования выбирают чаще всего?"
python main.py "удали все данные"
```

В консоли будут debug-логи tools (`[tool:...]`, `[trace:...]`) и финальный ответ по контракту.

## Контракт ответа

Агент обязан отвечать строго в формате:

```text
Status: success | error
Action: <описание действия>
Data: <результат>
Errors: <ошибка или ->
```

Источник требований: [`prompts.py`](prompts.py) (system prompt) и [`PROMPTS.md`](PROMPTS.md).

## Tools

| Tool | Тип | API / действие |
|------|-----|----------------|
| `get_questions` | HTTP | `GET /questions` |
| `submit_answers` | HTTP | `POST /answers` |
| `list_submissions` | HTTP | `GET /submissions` |
| `get_submission` | HTTP | `GET /submissions/{id}` |
| `get_stats` | HTTP | `GET /stats` |
| `query_database` | NL→SQL | read-only SELECT к `survey.db` |

Объявление HTTP-tools и реальный вызов API: [`tools_api.py`](tools_api.py).  
Debug print результата tool: функция `_debug` в том же файле.

## Пример запрос → tool

- Запрос: `покажи вопросы анкеты`
- Ожидаемый tool / метод: `get_questions` → `GET /questions`

Подробные прогоны: [`REPORT.md`](REPORT.md), сырые логи: [`e2e-logs/`](e2e-logs/).  
Промпты: [`PROMPTS.md`](PROMPTS.md).  
Соответствие критериям ДЗ: [`CRITERIA.md`](CRITERIA.md).

## Секреты

- Реальные ключи только в `.env` (не коммитится).
- В репозитории — `.env.example` с заглушкой `lm-studio`.
