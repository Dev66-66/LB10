from fastapi import APIRouter

from models.workout import WorkoutCreate, WorkoutResponse
from services.go_workout_service import GoWorkoutService
from services.grpc_workout_service import GrpcWorkoutService

router = APIRouter(prefix="/workouts", tags=["workouts"])

_go = GoWorkoutService()
_grpc = GrpcWorkoutService()


@router.get("", response_model=list[WorkoutResponse])
async def list_workouts():
    return await _go.get_all()


@router.post("", response_model=WorkoutResponse, status_code=201)
async def create_workout(workout: WorkoutCreate):
    return await _go.create(workout.model_dump())


@router.get("/{workout_id}/grpc", response_model=WorkoutResponse)
async def get_workout_grpc(workout_id: int):
    """Fetch a single workout via gRPC from the Go service."""
    return await _grpc.get_by_id(workout_id)
