# План реализации: Мини-анкета (Full-Stack, Docker)

## Цель

Создать full-stack приложение «Мини-анкета» на Go (backend) + HTML/JS (frontend),
упакованное в Docker-контейнер, с хранением ответов в памяти.

---

## Архитектура

```
OtusAI-hw1/
├── backend/
│   ├── main.go          # HTTP-сервер: GET /questions, POST /answers, раздача статики
│   └── go.mod
├── frontend/
│   └── index.html       # Одностраничное SPA (HTML + CSS + JS)
├── screenshots/         # Скриншоты работающего приложения
├── Dockerfile           # Multi-stage сборка: Go binary + статика
├── docker-compose.yml   # Удобный запуск одной командой
├── .gitignore
├── README.md            # Инструкция по запуску через Docker
├── PROMPTS.md           # Все использованные промты
└── PLAN.md              # Этот файл
```

---

## Шаги реализации

### Шаг 1 — Структура проекта
Создать папки `backend/`, `frontend/`, `screenshots/`.

### Шаг 2 — Backend (Go)
- Инициализировать Go-модуль (`go mod init`)
- Реализовать HTTP-сервер на стандартной библиотеке `net/http`
- **`GET /questions`** — возвращает JSON с 3–5 вопросами (текст, тип, варианты ответа)
- **`POST /answers`** — принимает JSON `{answers: [...]}`, сохраняет в slice в памяти, возвращает `{"status": "ok"}`
- Раздача `frontend/index.html` как статики на корневом маршруте `/`
- CORS-заголовки для локальной разработки

### Шаг 3 — Frontend (HTML + JS)
- Один файл `index.html` со встроенными стилями и скриптом
- При загрузке: `fetch("/questions")` → динамическая генерация формы
- Поддержка разных типов вопросов: текст, radio, select
- По сабмиту: `POST "/answers"` → показать «Спасибо за ответы!»
- Красивый минималистичный UI (CSS без фреймворков)

### Шаг 4 — Dockerfile (multi-stage)
- **Stage 1 (builder)**: образ `golang:1.26-alpine`, сборка бинарника
- **Stage 2 (runtime)**: минимальный образ `alpine:latest`, копирует бинарник + `frontend/`
- Итоговый образ ~20 МБ

### Шаг 5 — docker-compose.yml
- Один сервис `app`, порт `8080:8080`
- Команда запуска: `docker compose up --build`

### Шаг 6 — README.md
- Краткое описание проекта
- Инструкция запуска через Docker Compose
- Инструкция запуска без Docker (локально через `go run`)
- Описание API endpoints

### Шаг 7 — PROMPTS.md
- Записать все промты, использованные при генерации кода

### Шаг 8 — Тестирование и скриншоты
- Запустить `docker compose up --build`
- Открыть `http://localhost:8080`
- Заполнить и отправить анкету
- Сделать скриншоты: форма, сообщение «Спасибо»
- Проверить ответы через curl или логи

### Шаг 9 — Git
- Добавить `.gitignore` (бинарники, node_modules и т.д.)
- Финальный коммит

---

## Технологический стек

| Слой       | Технология                        |
|------------|-----------------------------------|
| Backend    | Go 1.26, stdlib `net/http`        |
| Frontend   | HTML5 + CSS3 + Vanilla JS         |
| Упаковка   | Docker (multi-stage), Compose v2  |
| Хранилище  | In-memory slice ([]Submission)    |
| VCS        | Git                               |
