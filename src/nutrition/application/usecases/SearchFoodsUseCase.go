package usecases

import (
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type SearchFoodsUseCase struct {
	Repo interfaces.FoodRepository
}

func NewSearchFoodsUseCase(repo interfaces.FoodRepository) *SearchFoodsUseCase {
	return &SearchFoodsUseCase{Repo: repo}
}

func (uc *SearchFoodsUseCase) Execute(name string, page int, pageSize int) ([]models.Food, int64, error) {
	return uc.Repo.SearchByName(name, page, pageSize)
}
