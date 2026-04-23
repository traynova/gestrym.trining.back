package repositories

import (
	"gestrym-training/src/common/models"
	"gorm.io/gorm"
)

type NutritionPlanRepositoryImpl struct {
	db *gorm.DB
}

func NewNutritionPlanRepositoryImpl(db *gorm.DB) *NutritionPlanRepositoryImpl {
	return &NutritionPlanRepositoryImpl{db: db}
}

func (r *NutritionPlanRepositoryImpl) Save(plan *models.NutritionPlan) error {
	return r.db.Save(plan).Error
}

func (r *NutritionPlanRepositoryImpl) GetByUserID(userID uint) ([]models.NutritionPlan, error) {
	var plans []models.NutritionPlan
	err := r.db.Where("user_id = ?", userID).Find(&plans).Error
	return plans, err
}

func (r *NutritionPlanRepositoryImpl) GetActiveByUserID(userID uint) (*models.NutritionPlan, error) {
	var plan models.NutritionPlan
	err := r.db.Where("user_id = ? AND is_active = ?", userID, true).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *NutritionPlanRepositoryImpl) DeactivateAllForUser(userID uint) error {
	return r.db.Model(&models.NutritionPlan{}).
		Where("user_id = ? AND is_active = ?", userID, true).
		Update("is_active", false).Error
}
