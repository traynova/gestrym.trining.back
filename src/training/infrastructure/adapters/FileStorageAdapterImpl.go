package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"path/filepath"

	"gestrym-training/src/training/domain/interfaces"
)

type FileStorageAdapterImpl struct {
	BaseURL string
	APIKey  string
}

func NewFileStorageAdapterImpl(baseURL, apiKey string) interfaces.FileStorageAdapter {
	return &FileStorageAdapterImpl{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

func (a *FileStorageAdapterImpl) UploadFromURL(url string, service string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch file from URL: status %d", resp.StatusCode)
	}

	// Get filename from URL, removing query parameters
	filename := "file"
	if base := filepath.Base(url); base != "" && base != "." && base != "/" {
		// Remove query params if any
		if idx := bytes.IndexAny([]byte(base), "?#"); idx != -1 {
			filename = string(base[:idx])
		} else {
			filename = base
		}
	}

	contentType := resp.Header.Get("Content-Type")

	// If content type is generic or missing, try to detect it or infer it
	if contentType == "" || contentType == "application/octet-stream" {
		// If it's from ExerciseDB image endpoint, we know it's a gif
		if (filepath.Base(url) == "image" || filepath.Base(url) == "image/") && (filepath.Ext(filename) == "" || filepath.Ext(filename) == ".gif") {
			contentType = "image/gif"
		} else {
			contentType = "image/gif"
		}
	}

	// Ensure filename has an appropriate extension if missing
	if filepath.Ext(filename) == "" {
		if contentType == "image/gif" {
			filename += ".gif"
		} else if contentType == "image/jpeg" {
			filename += ".jpg"
		} else if contentType == "image/png" {
			filename += ".png"
		}
	}

	log.Printf("Downloading file from URL: %s, detected contentType: %s, using filename: %s", url, contentType, filename)

	return a.UploadFromReader(resp.Body, filename, contentType, service)
}

func (a *FileStorageAdapterImpl) UploadFromReader(reader io.Reader, filename string, contentType string, service string) (string, error) {
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	errChan := make(chan error, 1)

	// Goroutine to write multipart data to the pipe
	go func() {
		defer pw.Close()
		defer writer.Close()

		// Add service field
		if err := writer.WriteField("service", service); err != nil {
			errChan <- fmt.Errorf("failed to write service field: %w", err)
			return
		}

		// Add files field with correct Content-Type
		h := make(textproto.MIMEHeader)
		h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="files"; filename="%s"`, filename))
		h.Set("Content-Type", contentType)
		part, err := writer.CreatePart(h)
		if err != nil {
			errChan <- fmt.Errorf("failed to create multipart part: %w", err)
			return
		}

		if _, err := io.Copy(part, reader); err != nil {
			errChan <- fmt.Errorf("failed to copy reader to part: %w", err)
			return
		}
	}()

	req, err := http.NewRequest("POST", a.BaseURL+"/internal/files/upload", pr)
	if err != nil {
		return "", fmt.Errorf("failed to create upload request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("X-API-Key", a.APIKey)
	req.Header.Set("accept", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute upload request: %w", err)
	}
	defer res.Body.Close()

	// Check for errors from the writer goroutine
	select {
	case err := <-errChan:
		return "", err
	default:
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(res.Body)
		return "", fmt.Errorf("storage service error (status %d): %s", res.StatusCode, string(respBody))
	}

	var result struct {
		CollectionID string `json:"collection_id"`
	}

	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("failed to decode storage service response: %w", err)
	}

	if result.CollectionID == "" {
		return "", fmt.Errorf("collectionId not found in storage response")
	}

	return result.CollectionID, nil
}
