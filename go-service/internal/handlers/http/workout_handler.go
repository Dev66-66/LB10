package http

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Dev66-66/LB10/go-service/internal/models"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

type WorkoutHandler struct {
	store *store.WorkoutStore
}

func NewWorkoutHandler(s *store.WorkoutStore) *WorkoutHandler {
	return &WorkoutHandler{store: s}
}

func (h *WorkoutHandler) RegisterRoutes(r *gin.Engine) {
	r.GET("/workouts", h.GetAll)
	r.POST("/workouts", h.Create)
	r.GET("/workouts/:id", h.GetByID)
	r.DELETE("/workouts/:id", h.Delete)
}

func (h *WorkoutHandler) GetAll(c *gin.Context) {
	c.JSON(http.StatusOK, h.store.GetAll())
}

type createWorkoutRequest struct {
	Name           string                  `json:"name"`
	Type           models.WorkoutType      `json:"type"`
	Duration       int                     `json:"duration"`
	Difficulty     models.WorkoutDifficulty `json:"difficulty"`
	CaloriesBurned int                     `json:"calories_burned"`
}

func (h *WorkoutHandler) Create(c *gin.Context) {
	var req createWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if strings.TrimSpace(req.Name) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name must not be empty"})
		return
	}
	if !req.Type.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type must be one of: cardio, strength, flexibility"})
		return
	}
	if !req.Difficulty.IsValid() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "difficulty must be one of: easy, medium, hard"})
		return
	}
	if req.Duration <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration must be positive"})
		return
	}
	if req.CaloriesBurned < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "calories_burned must not be negative"})
		return
	}

	w := h.store.Create(models.Workout{
		Name:           strings.TrimSpace(req.Name),
		Type:           req.Type,
		Duration:       req.Duration,
		Difficulty:     req.Difficulty,
		CaloriesBurned: req.CaloriesBurned,
	})

	c.JSON(http.StatusCreated, w)
}

func (h *WorkoutHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	w, err := h.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be an integer"})
		return
	}

	if err := h.store.Delete(id); err != nil {
		if errors.Is(err, store.ErrNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "workout not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.Status(http.StatusNoContent)
}
