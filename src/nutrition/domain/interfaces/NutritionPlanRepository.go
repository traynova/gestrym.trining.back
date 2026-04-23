package interfaces

import "gestrym-training/src/common/models"

type NutritionPlanRepository interface {
	Save(plan *models.NutritionPlan) error
	GetByUserID(userID uint) ([]models.NutritionPlan, error)
	GetActiveByUserID(userID uint) (*models.NutritionPlan, error)
	DeactivateAllForUser(userID uint) error
}
