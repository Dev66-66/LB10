package models

import "time"

type WorkoutType string
type WorkoutDifficulty string

const (
	TypeCardio      WorkoutType = "cardio"
	TypeStrength    WorkoutType = "strength"
	TypeFlexibility WorkoutType = "flexibility"
)

const (
	DifficultyEasy   WorkoutDifficulty = "easy"
	DifficultyMedium WorkoutDifficulty = "medium"
	DifficultyHard   WorkoutDifficulty = "hard"
)

type Workout struct {
	ID            int               `json:"id"`
	Name          string            `json:"name"`
	Type          WorkoutType       `json:"type"`
	Duration      int               `json:"duration"`
	Difficulty    WorkoutDifficulty `json:"difficulty"`
	CaloriesBurned int              `json:"calories_burned"`
	CreatedAt     time.Time         `json:"created_at"`
}

func (t WorkoutType) IsValid() bool {
	switch t {
	case TypeCardio, TypeStrength, TypeFlexibility:
		return true
	}
	return false
}

func (d WorkoutDifficulty) IsValid() bool {
	switch d {
	case DifficultyEasy, DifficultyMedium, DifficultyHard:
		return true
	}
	return false
}
