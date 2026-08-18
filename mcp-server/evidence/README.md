# Подтверждение вызовов из Cursor

Живой прогон Agent-чата (2026-08-18 21:49) зафиксирован в репозитории:

| Что | Где |
|-----|-----|
| Trace `CallToolRequest` + `[mcp]` | [`../logs/cursor-agent-stdio.log`](../logs/cursor-agent-stdio.log) |
| По запросам | [`../logs/01-get_questions.txt`](../logs/01-get_questions.txt) и соседние |
| JSON-ответы tools | [`../logs/cursor-results/`](../logs/cursor-results/) |
| Сводка | [`../REPORT.md`](../REPORT.md) |

Скриншоты UI не обязательны: Cursor пишет тот же вызов в MCP-лог (`CallToolRequest`).
