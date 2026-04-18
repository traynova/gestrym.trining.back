package models

type WorkoutSet struct {
	ID                uint    `gorm:"primaryKey" json:"id"`
	WorkoutExerciseID uint    `gorm:"not null;index" json:"workoutExerciseId"`
	Type              string  `gorm:"size:50;not null;default:'normal'" json:"type"` // warmup, normal, dropset
	Reps              int     `json:"reps"`
	Weight            float64 `json:"weight"`
	RestSeconds       int     `json:"restSeconds"`
}
