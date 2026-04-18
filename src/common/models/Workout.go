package models

import "time"

type Workout struct {
	ID        uint              `gorm:"primaryKey" json:"id"`
	UserID    uint              `gorm:"not null;index" json:"userId"`
	Name      string            `gorm:"size:255;not null" json:"name"`
	Exercises []WorkoutExercise `gorm:"foreignKey:WorkoutID" json:"exercises,omitempty"`
	CreatedAt time.Time         `json:"createdAt"`
	UpdatedAt time.Time         `json:"updatedAt"`
}
