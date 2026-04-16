package interfaces

import "gestrym-training/src/common/models"

// ExerciseDBAdapter defines how to interact with the external ExerciseDB API
type ExerciseDBAdapter interface {
	FetchAllExercises() ([]models.Exercise, error)
}
