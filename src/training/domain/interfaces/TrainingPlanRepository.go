package interfaces

import "gestrym-training/src/common/models"

// TrainingPlanRepository defines persistence operations for TrainingPlan entities.
type TrainingPlanRepository interface {
	// Create persists a new training plan and returns the saved entity (with ID populated).
	Create(plan *models.TrainingPlan) (*models.TrainingPlan, error)

	// FindByID retrieves a plan with its days (and each day's workout) preloaded.
	FindByID(id uint) (*models.TrainingPlan, error)

	// FindByUserID retrieves all plans assigned to a given user (AssignedTo = userID).
	FindByUserID(userID uint) ([]models.TrainingPlan, error)

	// AssignToUser updates the AssignedTo field of an existing plan.
	AssignToUser(planID uint, userID uint) error

	// FindTemplates returns all plans marked as IsTemplate = true.
	FindTemplates() ([]models.TrainingPlan, error)

	// FindLatestByUserID returns the most recently created plan for a user.
	FindLatestByUserID(userID uint) (*models.TrainingPlan, error)
}
