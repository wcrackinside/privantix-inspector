#!/usr/bin/env python3
"""
Fetch latest GitHub Release per unique repo in privantix_site/products.json
and write privantix_site/data/releases-cache.json for offline / fallback use.

Usage:
  python scripts/sync_releases_cache.py

Optional: set GITHUB_TOKEN in the environment for higher rate limits (5000 req/h).
"""

from __future__ import annotations

import json
import os
import sys
import urllib.error
import urllib.request
from datetime import datetime, timezone
from pathlib import Path

ROOT = Path(__file__).resolve().parents[1]
PRODUCTS = ROOT / "privantix_site" / "products.json"
OUT = ROOT / "privantix_site" / "data" / "releases-cache.json"
API = "https://api.github.com/repos/{owner}/{repo}/releases/latest"


def load_products() -> list:
    with PRODUCTS.open(encoding="utf-8") as f:
        data = json.load(f)
    return data if isinstance(data, list) else data.get("products", [])


def unique_repos(products: list) -> list[tuple[str, str]]:
    seen: set[str] = set()
    out: list[tuple[str, str]] = []
    for p in products:
        r = (p.get("repo") or "").strip()
        if "/" not in r:
            continue
        owner, name = r.split("/", 1)
        key = f"{owner}/{name}"
        if key not in seen:
            seen.add(key)
            out.append((owner, name))
    return out


def fetch_release(owner: str, repo: str) -> dict:
    url = API.format(owner=owner, repo=repo)
    headers = {
        "Accept": "application/vnd.github+json",
        "User-Agent": "privantix-sync-releases-cache",
    }
    token = os.environ.get("GITHUB_TOKEN") or os.environ.get("GH_TOKEN")
    if token:
        headers["Authorization"] = f"Bearer {token}"
    req = urllib.request.Request(url, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=60) as resp:
            body = resp.read().decode("utf-8")
            return {"release": json.loads(body), "error": None}
    except urllib.error.HTTPError as e:
        if e.code == 404:
            return {"release": None, "error": "404"}
        return {"release": None, "error": f"HTTP {e.code}"}
    except OSError as e:
        return {"release": None, "error": str(e)}


def main() -> int:
    if not PRODUCTS.is_file():
        print(f"Missing {PRODUCTS}", file=sys.stderr)
        return 1
    products = load_products()
    repos = unique_repos(products)
    payload: dict = {
        "fetched_at": datetime.now(timezone.utc).isoformat(),
        "repos": {},
    }
    for owner, repo in repos:
        key = f"{owner}/{repo}"
        print(f"Fetching {key} ...")
        payload["repos"][key] = fetch_release(owner, repo)

    OUT.parent.mkdir(parents=True, exist_ok=True)
    with OUT.open("w", encoding="utf-8") as f:
        json.dump(payload, f, indent=2, ensure_ascii=False)
        f.write("\n")
    print(f"Wrote {OUT}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
