package interfaces

import "io"

// FileStorageAdapter defines how to interact with the internal storage microservice
type FileStorageAdapter interface {
	UploadFromURL(url string, service string) (string, error)
	UploadFromReader(reader io.Reader, filename string, contentType string, service string) (string, error)
}
