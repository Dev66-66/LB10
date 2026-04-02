package store

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/Dev66-66/LB10/go-service/internal/models"
)

var ErrNotFound = errors.New("workout not found")
var ErrInvalidName = errors.New("name must not be empty")

type WorkoutStore struct {
	mu      sync.RWMutex
	items   map[int]models.Workout
	nextID  int
}

func NewWorkoutStore() *WorkoutStore {
	return &WorkoutStore{
		items:  make(map[int]models.Workout),
		nextID: 1,
	}
}

func (s *WorkoutStore) GetAll() []models.Workout {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]models.Workout, 0, len(s.items))
	for _, w := range s.items {
		result = append(result, w)
	}
	return result
}

func (s *WorkoutStore) GetByID(id int) (models.Workout, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	w, ok := s.items[id]
	if !ok {
		return models.Workout{}, ErrNotFound
	}
	return w, nil
}

func (s *WorkoutStore) Create(w models.Workout) (models.Workout, error) {
	if strings.TrimSpace(w.Name) == "" {
		return models.Workout{}, ErrInvalidName
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	w.ID = s.nextID
	w.CreatedAt = time.Now().UTC()
	s.items[w.ID] = w
	s.nextID++
	return w, nil
}

func (s *WorkoutStore) Delete(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.items[id]; !ok {
		return ErrNotFound
	}
	delete(s.items, id)
	return nil
}
