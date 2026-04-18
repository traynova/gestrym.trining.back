package extractor

import (
	"context"
	"fmt"
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type USDAExtractor struct {
	USDAAdapter interfaces.USDAAdapter
}

func NewUSDAExtractor(usda interfaces.USDAAdapter) *USDAExtractor {
	return &USDAExtractor{USDAAdapter: usda}
}

func (e *USDAExtractor) Extract(ctx context.Context, query string) ([]models.Food, error) {
	foods, err := e.USDAAdapter.SearchFoods(query)
	if err != nil {
		return nil, fmt.Errorf("failed to extract from USDA: %w", err)
	}
	return foods, nil
}
