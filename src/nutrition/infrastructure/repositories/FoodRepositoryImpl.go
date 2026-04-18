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

func (r *FoodRepositoryImpl) SearchByName(name string, page int, pageSize int) ([]models.Food, int64, error) {
	var foods []models.Food
	var total int64

	query := r.DB.Model(&models.Food{}).Preload("Category")
	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err = query.Offset(offset).Limit(pageSize).Find(&foods).Error
	return foods, total, err
}

func (r *FoodRepositoryImpl) FindByID(id uint) (*models.Food, error) {
	var food models.Food
	err := r.DB.Preload("Category").First(&food, id).Error
	return &food, err
}

func (r *FoodRepositoryImpl) FindByName(name string) (*models.Food, error) {
	var food models.Food
	err := r.DB.Where("name = ?", name).First(&food).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &food, nil
}

func (r *FoodRepositoryImpl) SaveFoods(foods []models.Food) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		for i := range foods {
			// Handle Category
			if foods[i].Category.Name != "" {
				var cat models.FoodCategory
				err := tx.Where("name = ?", foods[i].Category.Name).FirstOrCreate(&cat, models.FoodCategory{Name: foods[i].Category.Name}).Error
				if err != nil {
					return err
				}
				foods[i].CategoryID = cat.ID
				foods[i].Category = models.FoodCategory{} // Clear to avoid re-insertion
			}

			if err := tx.Create(&foods[i]).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
