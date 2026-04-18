package repositories

import (
	"errors"

	"gestrym-training/src/common/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ExerciseRepositoryImpl struct {
	DB *gorm.DB
}

func NewExerciseRepositoryImpl(db *gorm.DB) *ExerciseRepositoryImpl {
	return &ExerciseRepositoryImpl{DB: db}
}

func (r *ExerciseRepositoryImpl) BulkInsertExercises(exercises []models.Exercise) error {
	// Idempotent insertion using OnConflict (upsert/do nothing).
	return r.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ext_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"collection_id", "gif_url", "updated_at"}),
	}).CreateInBatches(exercises, 100).Error
}

func (r *ExerciseRepositoryImpl) FindByName(name string) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.DB.Where("name = ?", name).First(&exercise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil naturally if not found
		}
		return nil, err
	}
	return &exercise, nil
}

func (r *ExerciseRepositoryImpl) FindByExtID(extID string) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.DB.Where("ext_id = ?", extID).First(&exercise).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &exercise, nil
}

func (r *ExerciseRepositoryImpl) ListAll(bodyPart string, target string) ([]models.Exercise, error) {
	var exercises []models.Exercise
	query := r.DB.Model(&models.Exercise{})

	if bodyPart != "" {
		query = query.Where("body_part = ?", bodyPart)
	}
	if target != "" {
		query = query.Where("target = ?", target)
	}

	err := query.Find(&exercises).Error
	return exercises, err
}

func (r *ExerciseRepositoryImpl) FindByID(id uint) (*models.Exercise, error) {
	var exercise models.Exercise
	err := r.DB.First(&exercise, id).Error
	if err != nil {
		return nil, err
	}
	return &exercise, nil
}
