package repositories

import (
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/domain/interfaces"
	"gorm.io/gorm"
)

type FoodRepositoryImpl struct {
	DB *gorm.DB
}

func NewFoodRepositoryImpl(db *gorm.DB) interfaces.FoodRepository {
	return &FoodRepositoryImpl{DB: db}
}

func (r *FoodRepositoryImpl) SearchByName(name string) ([]models.Food, error) {
	var foods []models.Food
	err := r.DB.Preload("Category").Where("name ILIKE ?", "%"+name+"%").Find(&foods).Error
	return foods, err
}

func (r *FoodRepositoryImpl) FindByID(id uint) (*models.Food, error) {
	var food models.Food
	err := r.DB.Preload("Category").First(&food, id).Error
	if err != nil {
		return nil, err
	}
	return &food, nil
}

func (r *FoodRepositoryImpl) Create(food *models.Food) error {
	return r.DB.Create(food).Error
}
