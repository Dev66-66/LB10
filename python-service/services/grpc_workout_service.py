import os

import grpc
from grpc import aio
from fastapi import HTTPException

from proto import workout_pb2, workout_pb2_grpc

_GO_GRPC_HOST = os.getenv("GRPC_HOST", "localhost:50051")


class GrpcWorkoutService:
    """gRPC client for the Go workout service."""

    async def get_by_id(self, workout_id: int) -> dict:
        try:
            async with aio.insecure_channel(_GO_GRPC_HOST) as channel:
                stub = workout_pb2_grpc.WorkoutServiceStub(channel)
                response = await stub.GetWorkout(
                    workout_pb2.GetWorkoutRequest(id=workout_id)
                )
                w = response.workout
                return {
                    "id": w.id,
                    "name": w.name,
                    "type": w.type,
                    "duration": w.duration,
                    "difficulty": w.difficulty,
                    "calories_burned": w.calories_burned,
                    "created_at": w.created_at,
                }
        except grpc.aio.AioRpcError as e:
            if e.code() == grpc.StatusCode.NOT_FOUND:
                raise HTTPException(status_code=404, detail=e.details())
            raise HTTPException(status_code=503, detail="gRPC service unavailable")
