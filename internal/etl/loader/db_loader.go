package loader

import (
	"context"
	"fmt"
	"log"
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/application/utils"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type DBLoader struct {
	Repo           interfaces.FoodRepository
	StorageService interfaces.StorageService
}

func NewDBLoader(repo interfaces.FoodRepository, storage interfaces.StorageService) *DBLoader {
	return &DBLoader{
		Repo:           repo,
		StorageService: storage,
	}
}

func (l *DBLoader) Load(ctx context.Context, food models.Food, externalImageURL string) error {
	// 1. Deduplication (Final check)
	existing, _ := l.Repo.FindByName(food.Name)
	if existing != nil {
		return nil // Skip
	}

	// 2. Upload Image if URL provided
	if externalImageURL != "" {
		cleanName := utils.NormalizeFoodName(food.Name)
		fileName := fmt.Sprintf("%s.jpg", cleanName)
		
		collectionID, err := l.StorageService.UploadFromURL(ctx, externalImageURL, fileName)
		if err == nil {
			food.CollectionID = collectionID
			food.ImageURL = fmt.Sprintf("/storage/%s", collectionID)
		} else {
			log.Printf("[ETL ERROR] Failed to upload image for %s: %v", food.Name, err)
		}
	} else {
		log.Printf("[ETL WARNING] No external image URL provided for %s, skipping upload", food.Name)
	}

	// 3. Bulk Save (using our repository)
	return l.Repo.SaveFoods([]models.Food{food})
}
