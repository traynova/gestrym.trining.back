package interfaces

import "gestrym-training/src/common/models"

// TrainingDayRepository defines persistence operations for TrainingDay entities.
type TrainingDayRepository interface {
	// Create persists a new day entry within a training plan.
	Create(day *models.TrainingDay) (*models.TrainingDay, error)

	// FindByPlanID retrieves all days for a given plan, with Workout preloaded.
	FindByPlanID(planID uint) ([]models.TrainingDay, error)

	// FindByID retrieves a single training day by its ID.
	FindByID(dayID uint) (*models.TrainingDay, error)

	// UpdateCompletionStatus updates the IsCompleted field for a day.
	UpdateCompletionStatus(dayID uint, completed bool) error
}
