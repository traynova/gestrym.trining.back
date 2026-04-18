package interfaces

import "context"

type StorageService interface {
	UploadFromURL(ctx context.Context, imageURL string, fileName string) (string, error)
}
