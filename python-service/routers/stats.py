from collections import Counter

from fastapi import APIRouter

from services.go_workout_service import GoWorkoutService

router = APIRouter(prefix="/stats", tags=["stats"])

_go = GoWorkoutService()


@router.get("")
async def get_stats():
    """Aggregate workout statistics by type and difficulty."""
    workouts = await _go.get_all()
    by_type = Counter(w["type"] for w in workouts)
    by_difficulty = Counter(w["difficulty"] for w in workouts)
    return {
        "total": len(workouts),
        "by_type": dict(by_type),
        "by_difficulty": dict(by_difficulty),
    }
