package repositories

import (
	"gestrym-training/src/common/models"
	"gorm.io/gorm"
)

// AssignmentRepositoryImpl satisfies domain/interfaces.AssignmentRepository using GORM.
type AssignmentRepositoryImpl struct {
	DB *gorm.DB
}

func NewAssignmentRepositoryImpl(db *gorm.DB) *AssignmentRepositoryImpl {
	return &AssignmentRepositoryImpl{DB: db}
}

// Assign persists a new training plan assignment.
func (r *AssignmentRepositoryImpl) Assign(assignment *models.TrainingPlanAssignment) (*models.TrainingPlanAssignment, error) {
	if err := r.DB.Create(assignment).Error; err != nil {
		return nil, err
	}
	return assignment, nil
}

// FindByUserID returns all assignments for a user, preloading the TrainingPlan.
func (r *AssignmentRepositoryImpl) FindByUserID(userID uint) ([]models.TrainingPlanAssignment, error) {
	var assignments []models.TrainingPlanAssignment
	err := r.DB.
		Where("user_id = ?", userID).
		Preload("TrainingPlan").
		Find(&assignments).Error
	return assignments, err
}
