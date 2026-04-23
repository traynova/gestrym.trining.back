package repositories

import (
	"errors"

	"gestrym-training/src/common/models"
	"gorm.io/gorm"
)

// TrainingPlanRepositoryImpl satisfies domain/interfaces.TrainingPlanRepository using GORM.
type TrainingPlanRepositoryImpl struct {
	DB *gorm.DB
}

func NewTrainingPlanRepositoryImpl(db *gorm.DB) *TrainingPlanRepositoryImpl {
	return &TrainingPlanRepositoryImpl{DB: db}
}

// Create persists a new TrainingPlan and returns the saved record (ID populated by GORM).
func (r *TrainingPlanRepositoryImpl) Create(plan *models.TrainingPlan) (*models.TrainingPlan, error) {
	if err := r.DB.Create(plan).Error; err != nil {
		return nil, err
	}
	return plan, nil
}

// FindByID loads a plan with Days → Workout preloaded to avoid N+1 queries.
func (r *TrainingPlanRepositoryImpl) FindByID(id uint) (*models.TrainingPlan, error) {
	var plan models.TrainingPlan
	err := r.DB.
		Preload("Days").
		Preload("Days.Workout").
		Preload("Days.Workout.Exercises").
		Preload("Days.Workout.Exercises.Exercise").
		First(&plan, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// FindByUserID returns all plans where AssignedTo = userID, with days preloaded.
func (r *TrainingPlanRepositoryImpl) FindByUserID(userID uint) ([]models.TrainingPlan, error) {
	var plans []models.TrainingPlan
	err := r.DB.
		Where("assigned_to = ?", userID).
		Preload("Days").
		Preload("Days.Workout").
		Find(&plans).Error

	return plans, err
}

// AssignToUser sets the AssignedTo field on an existing plan.
func (r *TrainingPlanRepositoryImpl) AssignToUser(planID uint, userID uint) error {
	return r.DB.Model(&models.TrainingPlan{}).
		Where("id = ?", planID).
		Update("assigned_to", userID).Error
}

// FindTemplates returns plans flagged as reusable templates.
func (r *TrainingPlanRepositoryImpl) FindTemplates() ([]models.TrainingPlan, error) {
	var plans []models.TrainingPlan
	err := r.DB.Where("is_template = ?", true).Find(&plans).Error
	return plans, err
}

// FindLatestByUserID returns the most recently created plan for a user.
func (r *TrainingPlanRepositoryImpl) FindLatestByUserID(userID uint) (*models.TrainingPlan, error) {
	var plan models.TrainingPlan
	err := r.DB.
		Where("assigned_to = ?", userID).
		Preload("Days").
		Preload("Days.Workout").
		Order("created_at DESC").
		First(&plan).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}
