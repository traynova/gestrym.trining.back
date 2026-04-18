package transformer

import (
	"context"
	"fmt"
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/application/utils"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type FoodTransformer struct {
	ImageProvider interfaces.ImageProvider
}

func NewFoodTransformer(imageProvider interfaces.ImageProvider) *FoodTransformer {
	return &FoodTransformer{ImageProvider: imageProvider}
}

func (t *FoodTransformer) Transform(ctx context.Context, food models.Food) (models.Food, string, error) {
	// 1. Normalize
	cleanName := utils.NormalizeFoodName(food.Name)
	
	// 2. Search Image
	externalURL, err := t.ImageProvider.SearchImage(cleanName)
	if err != nil {
		return food, "", fmt.Errorf("image search failed for %s: %w", cleanName, err)
	}

	return food, externalURL, nil
}
