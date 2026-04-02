from typing import Annotated

from fastapi import APIRouter, Depends

from models.workout import WorkoutCreate, WorkoutResponse
from services.dependencies import get_go_service, get_grpc_service
from services.go_workout_service import GoWorkoutService
from services.grpc_workout_service import GrpcWorkoutService

router = APIRouter(prefix="/workouts", tags=["workouts"])

_Go = Annotated[GoWorkoutService, Depends(get_go_service)]
_Grpc = Annotated[GrpcWorkoutService, Depends(get_grpc_service)]


@router.get("", response_model=list[WorkoutResponse])
async def list_workouts(go: _Go):
    return await go.get_all()


@router.post("", response_model=WorkoutResponse, status_code=201)
async def create_workout(workout: WorkoutCreate, go: _Go):
    return await go.create(workout.model_dump())


@router.get("/{workout_id}/grpc", response_model=WorkoutResponse)
async def get_workout_grpc(workout_id: int, grpc: _Grpc):
    """Fetch a single workout via gRPC from the Go service."""
    return await grpc.get_by_id(workout_id)
