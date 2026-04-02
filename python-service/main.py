from contextlib import asynccontextmanager

from fastapi import FastAPI

from routers import auth, stats, workouts
from services.dependencies import _go_service


@asynccontextmanager
async def lifespan(app: FastAPI):
    yield
    await _go_service.aclose()


app = FastAPI(
    title="Workout Tracker — Python Service",
    description="FastAPI proxy and aggregator for the Go workout service.",
    version="1.0.0",
    lifespan=lifespan,
)

app.include_router(auth.router)
app.include_router(workouts.router)
app.include_router(stats.router)
