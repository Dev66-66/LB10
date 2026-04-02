package app

import (
	"os"

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

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev-secret"
	}

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Logger())

	// Public routes
	auth := handlers.NewAuthHandler(secret)
	r.POST("/auth/login", auth.Login)

	// Protected routes
	protected := r.Group("/")
	protected.Use(middleware.JWT(secret))
	handlers.NewWorkoutHandler(s).RegisterRoutes(protected)

	return &App{router: r}
}

func (a *App) Run(addr string) error {
	return a.router.Run(addr)
}
