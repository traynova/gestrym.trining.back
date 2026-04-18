package interfaces

import "gestrym-training/src/common/models"

type FoodRepository interface {
	SearchByName(name string) ([]models.Food, error)
	FindByID(id uint) (*models.Food, error)
	Create(food *models.Food) error
}
