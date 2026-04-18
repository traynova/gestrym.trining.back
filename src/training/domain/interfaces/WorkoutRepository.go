package interfaces

import "gestrym-training/src/common/models"

type WorkoutRepository interface {
	FindFullWorkoutByID(id uint) (*models.Workout, error)
	Create(workout *models.Workout) error
}
