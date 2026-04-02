package app

import (
	"github.com/gin-gonic/gin"

	handlers "github.com/Dev66-66/LB10/go-service/internal/handlers/http"
	"github.com/Dev66-66/LB10/go-service/internal/middleware"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

type App struct {
	router *gin.Engine
}

func New(s *store.WorkoutStore) *App {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	h := handlers.NewWorkoutHandler(s)
	h.RegisterRoutes(r)

	return &App{router: r}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
