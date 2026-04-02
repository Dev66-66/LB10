import httpx
from fastapi import HTTPException

_GO_BASE_URL = "http://localhost:8080"
_GO_USERNAME = "admin"
_GO_PASSWORD = "password123"


class GoWorkoutService:
    """Async HTTP client for the Go workout service."""

    def __init__(self) -> None:
        self._client = httpx.AsyncClient(base_url=_GO_BASE_URL, timeout=10.0)
        self._token: str | None = None

    async def _authenticate(self) -> str:
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
        self._token = resp.json()["token"]
        return self._token

    async def _get_token(self) -> str:
        if self._token is None:
            return await self._authenticate()
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
