package interfaces

import "gestrym-training/src/common/models"

type FoodRepository interface {
	SearchByName(name string, page int, pageSize int) ([]models.Food, int64, error)
	FindByID(id uint) (*models.Food, error)
	FindByName(name string) (*models.Food, error)
	SaveFoods(foods []models.Food) error
}
