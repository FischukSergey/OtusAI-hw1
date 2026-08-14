"""Сборка LangChain/LangGraph агента."""

from __future__ import annotations

import os

from dotenv import load_dotenv
from langchain_openai import ChatOpenAI
from langgraph.prebuilt import create_react_agent

from prompts import SYSTEM_PROMPT
from tools_api import API_TOOLS
from tools_sql import build_nl_sql_tool

load_dotenv()


def build_llm() -> ChatOpenAI:
    return ChatOpenAI(
        base_url=os.getenv("OPENAI_API_BASE", "http://127.0.0.1:1234/v1"),
        api_key=os.getenv("OPENAI_API_KEY", "lm-studio"),
        model=os.getenv("OPENAI_MODEL", "local-model"),
        temperature=0,
    )


def build_agent():
    llm = build_llm()
    tools = [*API_TOOLS, build_nl_sql_tool(llm)]
    return create_react_agent(llm, tools, prompt=SYSTEM_PROMPT)


def run_agent(user_query: str) -> str:
    agent = build_agent()
    print(f"[agent] user_query={user_query!r}", flush=True)
    result = agent.invoke(
        {"messages": [{"role": "user", "content": user_query}]},
        config={"recursion_limit": 20},
    )
    messages = result.get("messages", [])
    for msg in messages:
        kind = type(msg).__name__
        if hasattr(msg, "tool_calls") and msg.tool_calls:
            print(f"[trace:{kind}] tool_calls={msg.tool_calls}", flush=True)
        if kind == "ToolMessage":
            print(f"[trace:ToolMessage] name={getattr(msg, 'name', None)} content={msg.content}", flush=True)

    if not messages:
        return (
            "Status: error\n"
            "Action: none\n"
            "Data: -\n"
            "Errors: пустой ответ агента"
        )

    last = messages[-1]
    content = getattr(last, "content", str(last))
    if isinstance(content, list):
        content = "".join(
            part.get("text", "") if isinstance(part, dict) else str(part) for part in content
        )
    return str(content)
