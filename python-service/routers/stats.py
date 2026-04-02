from collections import Counter
from typing import Annotated

from fastapi import APIRouter, Depends

from services.dependencies import get_go_service
from services.go_workout_service import GoWorkoutService

router = APIRouter(prefix="/stats", tags=["stats"])

_Go = Annotated[GoWorkoutService, Depends(get_go_service)]


@router.get("")
async def get_stats(go: _Go):
    """Aggregate workout statistics by type and difficulty."""
    workouts = await go.get_all()
    by_type = Counter(w["type"] for w in workouts)
    by_difficulty = Counter(w["difficulty"] for w in workouts)
    return {
        "total": len(workouts),
        "by_type": dict(by_type),
        "by_difficulty": dict(by_difficulty),
    }
