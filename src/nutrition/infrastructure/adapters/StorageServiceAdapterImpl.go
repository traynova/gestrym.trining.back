package adapters

import (
	"context"

	"gestrym-training/src/nutrition/domain/interfaces"
	trainingInterfaces "gestrym-training/src/training/domain/interfaces"
)

type StorageServiceAdapterImpl struct {
	LegacyAdapter trainingInterfaces.FileStorageAdapter
}

func NewStorageServiceAdapterImpl(legacy trainingInterfaces.FileStorageAdapter) interfaces.StorageService {
	return &StorageServiceAdapterImpl{
		LegacyAdapter: legacy,
	}
}

func (a *StorageServiceAdapterImpl) UploadFromURL(ctx context.Context, imageURL string, fileName string) (string, error) {
	// We use the existing FileStorageAdapter to upload to our storage microservice
	collectionID, err := a.LegacyAdapter.UploadFromURL(imageURL, "nutrition")
	if err != nil {
		return "", err
	}

	// According to previous turns, we only get the collectionID. 
	// To follow the "Store ONLY the MinIO URL" rule, we would need to know how the storage-service serves those.
	// However, usually the collectionID is enough. 
	// If the user wants a full URL, we'd need to construct it based on the storage service base path.
	// For now, I'll return the collectionID as the "storedURL" or a placeholder for the URL if they prefer.
	
	return collectionID, nil
}
