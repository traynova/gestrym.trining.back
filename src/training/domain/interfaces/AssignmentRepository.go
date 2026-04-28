package interfaces

import "gestrym-training/src/common/models"

// AssignmentRepository defines persistence operations for TrainingPlanAssignments.
type AssignmentRepository interface {
	// Assign persists a new training plan assignment record.
	Assign(assignment *models.TrainingPlanAssignment) (*models.TrainingPlanAssignment, error)

	// FindByUserID retrieves all assignments for a specific user.
	FindByUserID(userID uint) ([]models.TrainingPlanAssignment, error)
}
