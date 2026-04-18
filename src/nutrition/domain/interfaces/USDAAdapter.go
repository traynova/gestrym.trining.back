package interfaces

import "gestrym-training/src/common/models"

type USDAAdapter interface {
	SearchFoods(query string) ([]models.Food, error)
}
