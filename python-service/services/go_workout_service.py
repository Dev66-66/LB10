import asyncio
import os

import httpx
from fastapi import HTTPException

_GO_BASE_URL = os.getenv("GO_SERVICE_URL", "http://localhost:8080")
_GO_USERNAME = os.getenv("GO_USERNAME", "admin")
_GO_PASSWORD = os.getenv("GO_PASSWORD", "password123")


class GoWorkoutService:
    """Async HTTP client for the Go workout service."""

    def __init__(self) -> None:
        self._client = httpx.AsyncClient(base_url=_GO_BASE_URL, timeout=10.0)
        self._token: str | None = None
        self._token_lock = asyncio.Lock()

    async def _authenticate(self) -> str:
        """Fetch a fresh token from Go /auth/login. Must be called under _token_lock."""
        try:
            resp = await self._client.post(
                "/auth/login",
                json={"username": _GO_USERNAME, "password": _GO_PASSWORD},
            )
            resp.raise_for_status()
        except httpx.ConnectError:
            raise HTTPException(status_code=503, detail="Go service unavailable")
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=e.response.text,
            )
        return resp.json()["token"]

    async def _get_token(self) -> str:
        if self._token is not None:
            return self._token
        async with self._token_lock:
            # Double-checked: another coroutine may have filled it while we waited.
            if self._token is None:
                self._token = await self._authenticate()
            return self._token

    async def _request(self, method: str, path: str, **kwargs) -> dict | list:
        token = await self._get_token()
        headers = {"Authorization": f"Bearer {token}"}
        try:
            resp = await self._client.request(method, path, headers=headers, **kwargs)
            resp.raise_for_status()
            return resp.json()
        except httpx.ConnectError:
            raise HTTPException(status_code=503, detail="Go service unavailable")
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=e.response.text,
            )

    async def get_all(self) -> list:
        return await self._request("GET", "/workouts")

    async def create(self, data: dict) -> dict:
        return await self._request("POST", "/workouts", json=data)

    async def fetch_token(self, username: str, password: str) -> str:
        """Fetch a JWT from the Go service for the given credentials."""
        try:
            resp = await self._client.post(
                "/auth/login",
                json={"username": username, "password": password},
            )
            resp.raise_for_status()
            return resp.json()["token"]
        except httpx.ConnectError:
            raise HTTPException(status_code=503, detail="Go service unavailable")
        except httpx.HTTPStatusError as e:
            raise HTTPException(
                status_code=e.response.status_code,
                detail=e.response.text,
            )

    async def aclose(self) -> None:
        await self._client.aclose()
