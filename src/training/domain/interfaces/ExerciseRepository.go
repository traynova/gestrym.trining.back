package interfaces

import "gestrym-training/src/common/models"

// ExerciseRepository defines the operations to persist exercises in our DB
type ExerciseRepository interface {
	BulkInsertExercises(exercises []models.Exercise) error
	FindByName(name string) (*models.Exercise, error)
	FindByExtID(extID string) (*models.Exercise, error)
	ListAll(bodyPart string, target string) ([]models.Exercise, error)
	FindByID(id uint) (*models.Exercise, error)
}
