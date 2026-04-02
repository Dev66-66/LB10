from typing import Annotated

from fastapi import APIRouter, Depends
from pydantic import BaseModel

from services.dependencies import get_go_service
from services.go_workout_service import GoWorkoutService

router = APIRouter(prefix="/auth", tags=["auth"])

_Go = Annotated[GoWorkoutService, Depends(get_go_service)]


class TokenRequest(BaseModel):
    username: str
    password: str


class TokenResponse(BaseModel):
    token: str


@router.post("/token", response_model=TokenResponse)
async def get_token(body: TokenRequest, go: _Go):
    """Proxy: obtain a JWT from the Go service and return it to the client."""
    token = await go.fetch_token(body.username, body.password)
    return TokenResponse(token=token)
