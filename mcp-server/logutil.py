"""Debug-лог MCP-инструментов. Только stderr: stdout занят транспортом stdio."""

from __future__ import annotations

import json
import sys
from datetime import datetime, timezone
from typing import Any


def log_tool(
    name: str,
    params: dict[str, Any],
    status: str,
    error: str | None = None,
) -> None:
    """Пишет имя tool, входные параметры (без секретов) и status=success|error."""
    payload: dict[str, Any] = {
        "ts": datetime.now(timezone.utc).isoformat(),
        "tool": name,
        "params": params,
        "status": status,
    }
    if error:
        payload["error"] = error
    print(f"[mcp] {json.dumps(payload, ensure_ascii=False)}", file=sys.stderr, flush=True)
