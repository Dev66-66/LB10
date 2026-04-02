from fastapi import APIRouter
from pydantic import BaseModel

from services.go_workout_service import GoWorkoutService

router = APIRouter(prefix="/auth", tags=["auth"])

_go = GoWorkoutService()


class TokenRequest(BaseModel):
    username: str
    password: str


class TokenResponse(BaseModel):
    token: str


@router.post("/token", response_model=TokenResponse)
async def get_token(body: TokenRequest):
    """Proxy: obtain a JWT from the Go service and return it to the client."""
    token = await _go.fetch_token(body.username, body.password)
    return TokenResponse(token=token)
