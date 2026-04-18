package usecases

import (
	"log"

	"gestrym-training/src/training/domain/interfaces"
)

// ImportExercisesUseCase orchestrates the fetching and storing of exercises
type ImportExercisesUseCase struct {
	Adapter        interfaces.ExerciseDBAdapter
	StorageAdapter interfaces.FileStorageAdapter
	Repo           interfaces.ExerciseRepository
}

func NewImportExercisesUseCase(adapter interfaces.ExerciseDBAdapter, storage interfaces.FileStorageAdapter, repo interfaces.ExerciseRepository) *ImportExercisesUseCase {
	return &ImportExercisesUseCase{
		Adapter:        adapter,
		StorageAdapter: storage,
		Repo:           repo,
	}
}

// Execute pulls data from the external API and inserts/updates the local DB idempotently
func (uc *ImportExercisesUseCase) Execute() error {
	log.Println("Starting exercise import from external API...")

	exercises, err := uc.Adapter.FetchAllExercises()
	if err != nil {
		return err
	}

	for i := range exercises {
		if exercises[i].GifURL != "" {
			log.Printf("Uploading gif for exercise: %s", exercises[i].Name)
			collectionID, err := uc.StorageAdapter.UploadFromURL(exercises[i].GifURL, "training")
			if err != nil {
				log.Printf("Warning: failed to upload gif for %s: %v", exercises[i].Name, err)
				continue
			}
			exercises[i].CollectionID = collectionID
		}
	}

	log.Printf("Fetched %d exercises. Saving to database... (ignoring existing to maintain idempotency)", len(exercises))

	// In bulk insertion mode, GORM can handle idempotency perfectly if unique indexes are defined (OnConflict constraints).
	err = uc.Repo.BulkInsertExercises(exercises)
	if err != nil {
		return err
	}

	log.Println("Import completed successfully")
	return nil
}
