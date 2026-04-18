package usecases

import (
	"context"
	"fmt"
	"log"
	"sync"

	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/application/utils"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type ImportFoodsWithImagesUseCase struct {
	Repo           interfaces.FoodRepository
	USDAAdapter    interfaces.USDAAdapter
	ImageProvider  interfaces.ImageProvider
	StorageService interfaces.StorageService
}

func NewImportFoodsWithImagesUseCase(
	repo interfaces.FoodRepository,
	usda interfaces.USDAAdapter,
	imageProvider interfaces.ImageProvider,
	storage interfaces.StorageService,
) *ImportFoodsWithImagesUseCase {
	return &ImportFoodsWithImagesUseCase{
		Repo:           repo,
		USDAAdapter:    usda,
		ImageProvider:  imageProvider,
		StorageService: storage,
	}
}

func (uc *ImportFoodsWithImagesUseCase) Execute(ctx context.Context) error {
	seedQuery := []string{"chicken", "beef", "rice", "egg", "milk", "fish", "potato", "banana"}
	
	// Channels for worker pool
	foodChan := make(chan models.Food, 20)
	errChan := make(chan error, 20)
	var wg sync.WaitGroup

	// Start Workers
	workerCount := 5 // Processing 5 foods concurrently
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for food := range foodChan {
				if err := uc.processFood(ctx, food); err != nil {
					log.Printf("[ERROR] Failed to process food %s: %v", food.Name, err)
					errChan <- err
				}
			}
		}()
	}

	// Fetch and send foods to channel
	for _, query := range seedQuery {
		log.Printf("[INFO] Fetching metadata for query: %s", query)
		foods, err := uc.USDAAdapter.SearchFoods(query)
		if err != nil {
			log.Printf("[WARNING] USDA search failed for %s: %v", query, err)
			continue
		}

		for _, food := range foods {
			// Check if already in queue or processed (basic deduplication)
			existing, _ := uc.Repo.FindByName(food.Name)
			if existing != nil {
				continue
			}
			foodChan <- food
		}
	}

	close(foodChan)
	wg.Wait()
	close(errChan)

	log.Printf("[SUCCESS] Food import completed")
	return nil
}

func (uc *ImportFoodsWithImagesUseCase) processFood(ctx context.Context, food models.Food) error {
	// 1. Normalize name
	cleanName := utils.NormalizeFoodName(food.Name)
	
	// 2. Search Image
	externalImageURL, err := uc.ImageProvider.SearchImage(cleanName)
	if err != nil {
		return fmt.Errorf("image search failed for %s: %w", cleanName, err)
	}

	// 3. Upload to Storage (Streaming)
	fileName := fmt.Sprintf("%s.jpg", cleanName)
	collectionID, err := uc.StorageService.UploadFromURL(ctx, externalImageURL, fileName)
	if err != nil {
		return fmt.Errorf("storage upload failed for %s: %w", cleanName, err)
	}

	food.CollectionID = collectionID
	food.ImageURL = fmt.Sprintf("/storage/%s", collectionID)

	// 4. Save to DB
	if err := uc.Repo.SaveFoods([]models.Food{food}); err != nil {
		return fmt.Errorf("db save failed for %s: %w", food.Name, err)
	}

	log.Printf("[OK] Processed and saved: %s", food.Name)
	return nil
}
