package repositories

import (
	"gestrym-training/src/common/models"
	"gorm.io/gorm"
)

// TrainingDayRepositoryImpl satisfies domain/interfaces.TrainingDayRepository using GORM.
type TrainingDayRepositoryImpl struct {
	DB *gorm.DB
}

func NewTrainingDayRepositoryImpl(db *gorm.DB) *TrainingDayRepositoryImpl {
	return &TrainingDayRepositoryImpl{DB: db}
}

// Create persists a single TrainingDay and returns the saved record.
func (r *TrainingDayRepositoryImpl) Create(day *models.TrainingDay) (*models.TrainingDay, error) {
	if err := r.DB.Create(day).Error; err != nil {
		return nil, err
	}
	return day, nil
}

// FindByPlanID retrieves all days for a plan with their Workouts preloaded.
func (r *TrainingDayRepositoryImpl) FindByPlanID(planID uint) ([]models.TrainingDay, error) {
	var days []models.TrainingDay
	err := r.DB.
		Where("training_plan_id = ?", planID).
		Preload("Workout").
		Preload("Workout.Exercises").
		Preload("Workout.Exercises.Exercise").
		Order("day_number ASC").
		Find(&days).Error
	return days, err
}

// FindByID retrieves a single day by its ID.
func (r *TrainingDayRepositoryImpl) FindByID(dayID uint) (*models.TrainingDay, error) {
	var day models.TrainingDay
	err := r.DB.First(&day, dayID).Error
	if err != nil {
		return nil, err
	}
	return &day, nil
}

// UpdateCompletionStatus toggles the IsCompleted flag for a day.
func (r *TrainingDayRepositoryImpl) UpdateCompletionStatus(dayID uint, completed bool) error {
	return r.DB.Model(&models.TrainingDay{}).Where("id = ?", dayID).Update("is_completed", completed).Error
}
