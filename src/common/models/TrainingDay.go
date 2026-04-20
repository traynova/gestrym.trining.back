package models

import "time"

// TrainingDay represents a single day within a TrainingPlan, linked to an existing Workout.
type TrainingDay struct {
	ID             uint      `gorm:"primaryKey"              json:"id"`
	TrainingPlanID uint      `gorm:"not null;index"          json:"trainingPlanId"`
	DayNumber      int       `gorm:"not null"                json:"dayNumber"` // 1...N (within the plan duration)
	WorkoutID      uint      `gorm:"not null;index"          json:"workoutId"`
	Notes          string    `gorm:"type:text"               json:"notes"`
	IsCompleted    bool      `gorm:"default:false"           json:"isCompleted"`
	Workout        Workout   `gorm:"foreignKey:WorkoutID"    json:"workout,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}
