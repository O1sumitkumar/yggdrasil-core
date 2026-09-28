#!/usr/bin/env python3
"""Call the local OpenAI-compatible API. Python 3 standard library only."""

import json
import os
import urllib.request

BASE = os.environ.get("YGG_BASE_URL", "http://127.0.0.1:7331")
MODEL = os.environ.get("YGG_MODEL", "profile:general-assistant")
API_KEY = os.environ.get("YGG_API_KEY", "")


def headers():
    h = {"Content-Type": "application/json"}
    if API_KEY:
        h["Authorization"] = "Bearer " + API_KEY
    return h


def get_models():
    req = urllib.request.Request(BASE + "/v1/models", headers=headers(), method="GET")
    with urllib.request.urlopen(req) as resp:
        return json.load(resp)


def chat(content):
    payload = {
        "model": MODEL,
        "messages": [{"role": "user", "content": content}],
    }
    data = json.dumps(payload).encode()
    req = urllib.request.Request(
        BASE + "/v1/chat/completions", data=data, headers=headers(), method="POST"
    )
    with urllib.request.urlopen(req) as resp:
        return json.load(resp)


if __name__ == "__main__":
    print(json.dumps(get_models(), indent=2))
    print(json.dumps(chat("Hello"), indent=2))
