package repositories

import (
	"gestrym-training/src/common/models"
	"gestrym-training/src/training/domain/interfaces"
	"gorm.io/gorm"
)

type WorkoutRepositoryImpl struct {
	DB *gorm.DB
}

func NewWorkoutRepositoryImpl(db *gorm.DB) interfaces.WorkoutRepository {
	return &WorkoutRepositoryImpl{DB: db}
}

func (r *WorkoutRepositoryImpl) FindFullWorkoutByID(id uint) (*models.Workout, error) {
	var workout models.Workout
	err := r.DB.Preload("Exercises.Exercise").Preload("Exercises.Sets").First(&workout, id).Error
	if err != nil {
		return nil, err
	}
	return &workout, nil
}

func (r *WorkoutRepositoryImpl) Create(workout *models.Workout) error {
	return r.DB.Create(workout).Error
}
