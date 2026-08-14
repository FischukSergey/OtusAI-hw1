# Мини-анкета + LangChain-агент

Full-stack приложение «Мини-анкета» на Go + HTML/JS и NL-агент на LangChain.  
Домашнее задание курса OtusAI.

## Описание

- Backend на Go предоставляет REST API и хранит данные в **SQLite**
- Frontend — одна HTML-страница с динамической формой
- Папка [`agent/`](agent/) — LangChain-агент: естественный язык → API tools / NL→SQL

## Быстрый старт (Docker)

**Требования:** [Docker](https://docs.docker.com/get-docker/) и Docker Compose v2

```bash
git clone https://github.com/sergeymac/otusai-hw1.git
cd otusai-hw1
docker compose up --build
open http://localhost:8090
```

Остановить: `Ctrl+C`, затем `docker compose down`

Данные SQLite сохраняются в volume `survey-data`.  
Снаружи контейнер доступен на порту **8090** (внутри сервиса по-прежнему `8080`).

## Запуск без Docker (локально)

**Требования:** Go 1.25+

```bash
cd backend
PORT=8090 go run .
# Сервер: http://localhost:8090
# БД: ./data/survey.db (создаётся автоматически)
```

Изменить порт / путь к БД:

```bash
PORT=8090 DB_PATH=./data/survey.db go run .
```

## API

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/questions` | Список вопросов анкеты |
| POST | `/answers` | Сохранить ответы (в SQLite) |
| GET | `/submissions` | Список заполненных анкет |
| GET | `/submissions/{id}` | Одна анкета |
| GET | `/stats` | Агрегированная статистика |

Базовый URL в примерах: `http://localhost:8090`.

### GET /questions

```bash
curl http://localhost:8090/questions
```

### POST /answers

```bash
curl -X POST http://localhost:8090/answers \
  -H "Content-Type: application/json" \
  -d '{"answers": [{"question_id": 1, "value": "Иван"}, {"question_id": 2, "value": "Go"}, {"question_id": 3, "value": "1–3 года"}, {"question_id": 4, "value": "AI"}, {"question_id": 5, "value": "Соцсети"}]}'
```

Ответ: `{"status":"ok","id":1}`

### GET /stats

```bash
curl http://localhost:8090/stats
```

## LangChain-агент

Инструкция по запуску (LM Studio + CLI): **[`agent/README.md`](agent/README.md)**.

Кратко:

```bash
# 1) API уже запущен на :8090
# 2) LM Studio Local Server на :1234
cd agent
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
cp .env.example .env
python main.py "покажи вопросы анкеты"
```

### Контракт ответа агента

```text
Status: success | error
Action: <описание>
Data: <результат>
Errors: <ошибка или ->
```

Подробности: [`agent/prompts.py`](agent/prompts.py), [`agent/PROMPTS.md`](agent/PROMPTS.md).

### Подтверждение критериев ДЗ

Сводная таблица ссылок (файл/строки tools, debug-лог, примеры): [`agent/CRITERIA.md`](agent/CRITERIA.md).  
Проверочные запросы: [`agent/REPORT.md`](agent/REPORT.md).

## Структура проекта

```
.
├── backend/
│   ├── main.go          # HTTP-сервер
│   ├── db.go            # SQLite store
│   └── go.mod
├── frontend/
│   └── index.html
├── agent/               # LangChain-агент (HW)
│   ├── main.py
│   ├── agent.py
│   ├── tools_api.py
│   ├── tools_sql.py
│   ├── prompts.py
│   ├── REPORT.md
│   ├── PROMPTS.md
│   └── README.md
├── Dockerfile
├── docker-compose.yml
└── README.md
```

## Скриншоты

| Форма анкеты | Сообщение «Спасибо» |
|---|---|
| ![Форма](screenshots/form.png) | ![Спасибо](screenshots/thank-you.png) |
