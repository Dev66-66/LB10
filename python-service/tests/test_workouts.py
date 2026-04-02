import pytest
import httpx
import respx
from httpx import ASGITransport, AsyncClient

from main import app

GO_BASE = "http://localhost:8080"
_TOKEN = "mock-jwt-token"
_AUTH_RESP = {"token": _TOKEN}
_WORKOUT = {
    "id": 1,
    "name": "Morning Run",
    "type": "cardio",
    "duration": 30,
    "difficulty": "easy",
    "calories_burned": 200,
    "created_at": "2024-01-01T10:00:00Z",
}


@pytest.fixture(autouse=True)
def reset_tokens():
    """Reset cached JWT tokens on all module-level service instances before each test."""
    import routers.workouts as wk
    import routers.stats as st
    import routers.auth as au

    wk._go._token = None
    st._go._token = None
    au._go._token = None
    yield


@pytest.fixture
async def client():
    async with AsyncClient(
        transport=ASGITransport(app=app), base_url="http://test"
    ) as ac:
        yield ac


async def test_create_workout_valid(client):
    async with respx.mock() as mock:
        mock.post(f"{GO_BASE}/auth/login").mock(
            return_value=httpx.Response(200, json=_AUTH_RESP)
        )
        mock.post(f"{GO_BASE}/workouts").mock(
            return_value=httpx.Response(201, json=_WORKOUT)
        )

        resp = await client.post(
            "/workouts",
            json={
                "name": "Morning Run",
                "type": "cardio",
                "duration": 30,
                "difficulty": "easy",
                "calories_burned": 200,
            },
        )

    assert resp.status_code == 201
    data = resp.json()
    assert data["name"] == "Morning Run"
    assert data["type"] == "cardio"
    assert data["id"] == 1


async def test_create_workout_empty_name(client):
    """Pydantic validator must reject empty/whitespace names before any HTTP call is made."""
    async with respx.mock(assert_all_called=False) as mock:
        resp = await client.post(
            "/workouts",
            json={
                "name": "   ",
                "type": "cardio",
                "duration": 30,
                "difficulty": "easy",
                "calories_burned": 200,
            },
        )

    assert resp.status_code == 422
    assert mock.calls.call_count == 0


async def test_get_workouts(client):
    async with respx.mock() as mock:
        mock.post(f"{GO_BASE}/auth/login").mock(
            return_value=httpx.Response(200, json=_AUTH_RESP)
        )
        mock.get(f"{GO_BASE}/workouts").mock(
            return_value=httpx.Response(200, json=[_WORKOUT])
        )

        resp = await client.get("/workouts")

    assert resp.status_code == 200
    data = resp.json()
    assert isinstance(data, list)
    assert len(data) == 1
    assert data[0]["name"] == "Morning Run"
    assert data[0]["type"] == "cardio"


async def test_go_service_unavailable(client):
    """ConnectError from Go service must surface as HTTP 503."""
    async with respx.mock() as mock:
        mock.post(f"{GO_BASE}/auth/login").mock(
            side_effect=httpx.ConnectError("Connection refused")
        )

        resp = await client.get("/workouts")

    assert resp.status_code == 503
    assert "unavailable" in resp.json().get("detail", "").lower()


async def test_stats_aggregation(client):
    """GET /stats must correctly count workouts by type and difficulty."""
    workouts = [
        {**_WORKOUT, "id": 1, "type": "cardio", "difficulty": "easy"},
        {**_WORKOUT, "id": 2, "type": "cardio", "difficulty": "hard"},
        {**_WORKOUT, "id": 3, "type": "strength", "difficulty": "medium"},
    ]

    async with respx.mock() as mock:
        mock.post(f"{GO_BASE}/auth/login").mock(
            return_value=httpx.Response(200, json=_AUTH_RESP)
        )
        mock.get(f"{GO_BASE}/workouts").mock(
            return_value=httpx.Response(200, json=workouts)
        )

        resp = await client.get("/stats")

    assert resp.status_code == 200
    data = resp.json()
    assert data["total"] == 3
    assert data["by_type"]["cardio"] == 2
    assert data["by_type"]["strength"] == 1
    assert data["by_difficulty"]["easy"] == 1
    assert data["by_difficulty"]["hard"] == 1
    assert data["by_difficulty"]["medium"] == 1


async def test_jwt_forwarded(client):
    """GoWorkoutService must attach 'Authorization: Bearer <token>' to requests to Go."""
    async with respx.mock() as mock:
        mock.post(f"{GO_BASE}/auth/login").mock(
            return_value=httpx.Response(200, json=_AUTH_RESP)
        )
        workouts_route = mock.get(f"{GO_BASE}/workouts").mock(
            return_value=httpx.Response(200, json=[_WORKOUT])
        )

        await client.get("/workouts")

    assert workouts_route.called
    sent_auth = workouts_route.calls[0].request.headers.get("authorization", "")
    assert sent_auth == f"Bearer {_TOKEN}"
