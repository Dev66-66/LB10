package store_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Dev66-66/LB10/go-service/internal/models"
	"github.com/Dev66-66/LB10/go-service/internal/store"
)

func TestGetAll_EmptyReturnsSliceNotNil(t *testing.T) {
	s := store.NewWorkoutStore()

	result := s.GetAll()

	assert.NotNil(t, result)
	assert.Len(t, result, 0)
}

func TestCreate_ValidWorkout(t *testing.T) {
	s := store.NewWorkoutStore()
	input := models.Workout{
		Name:           "Morning Run",
		Type:           models.TypeCardio,
		Duration:       30,
		Difficulty:     models.DifficultyEasy,
		CaloriesBurned: 250,
	}

	w, err := s.Create(input)

	require.NoError(t, err)
	assert.Equal(t, 1, w.ID)
	assert.Equal(t, "Morning Run", w.Name)
	assert.Equal(t, models.TypeCardio, w.Type)
	assert.Equal(t, 30, w.Duration)
	assert.Equal(t, models.DifficultyEasy, w.Difficulty)
	assert.Equal(t, 250, w.CaloriesBurned)
	assert.False(t, w.CreatedAt.IsZero())
}

func TestCreate_EmptyNameRejected(t *testing.T) {
	s := store.NewWorkoutStore()
	input := models.Workout{
		Name:       "",
		Type:       models.TypeStrength,
		Duration:   45,
		Difficulty: models.DifficultyMedium,
	}

	_, err := s.Create(input)

	require.ErrorIs(t, err, store.ErrInvalidName)
	assert.Len(t, s.GetAll(), 0)
}

func TestCreate_WhitespaceNameRejected(t *testing.T) {
	s := store.NewWorkoutStore()
	input := models.Workout{
		Name:       "   ",
		Type:       models.TypeStrength,
		Duration:   45,
		Difficulty: models.DifficultyMedium,
	}

	_, err := s.Create(input)

	require.ErrorIs(t, err, store.ErrInvalidName)
	assert.Len(t, s.GetAll(), 0)
}

func TestGetByID_NotFound(t *testing.T) {
	s := store.NewWorkoutStore()

	_, err := s.GetByID(999)

	require.ErrorIs(t, err, store.ErrNotFound)
}
