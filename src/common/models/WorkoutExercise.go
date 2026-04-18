package models

type WorkoutExercise struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	WorkoutID  uint         `gorm:"not null;index" json:"workoutId"`
	ExerciseID uint         `gorm:"not null" json:"exerciseId"`
	Exercise   Exercise     `gorm:"foreignKey:ExerciseID" json:"exercise,omitempty"`
	Order      int          `gorm:"not null;default:0" json:"order"`
	Sets       []WorkoutSet `gorm:"foreignKey:WorkoutExerciseID" json:"sets,omitempty"`
}
