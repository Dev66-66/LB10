from fastapi import FastAPI

from routers import auth, stats, workouts

app = FastAPI(
    title="Workout Tracker — Python Service",
    description="FastAPI proxy and aggregator for the Go workout service.",
    version="1.0.0",
)

app.include_router(auth.router)
app.include_router(workouts.router)
app.include_router(stats.router)
