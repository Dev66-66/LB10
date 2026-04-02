from typing import Literal

from pydantic import BaseModel, field_validator


class WorkoutCreate(BaseModel):
    name: str
    type: Literal["cardio", "strength", "flexibility"]
    duration: int
    difficulty: Literal["easy", "medium", "hard"]
    calories_burned: int

    @field_validator("name")
    @classmethod
    def name_not_empty(cls, v: str) -> str:
        if not v.strip():
            raise ValueError("name must not be empty or whitespace")
        return v.strip()


class WorkoutResponse(BaseModel):
    id: int
    name: str
    type: str
    duration: int
    difficulty: str
    calories_burned: int
    created_at: str
