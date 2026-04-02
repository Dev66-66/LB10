"""
Centralised service singletons and FastAPI dependency providers.
All routers import from here so the application shares exactly one
GoWorkoutService instance (one httpx connection pool, one token cache).
"""
from services.go_workout_service import GoWorkoutService
from services.grpc_workout_service import GrpcWorkoutService

_go_service = GoWorkoutService()
_grpc_service = GrpcWorkoutService()


def get_go_service() -> GoWorkoutService:
    return _go_service


def get_grpc_service() -> GrpcWorkoutService:
    return _grpc_service
