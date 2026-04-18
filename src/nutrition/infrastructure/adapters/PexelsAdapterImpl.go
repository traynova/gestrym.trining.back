package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type PexelsAdapterImpl struct {
	BaseURL string
	APIKey  string
}

func NewPexelsAdapterImpl(apiKey string) *PexelsAdapterImpl {
	return &PexelsAdapterImpl{
		BaseURL: "https://api.pexels.com/v1",
		APIKey:  apiKey,
	}
}

type pexelsResponse struct {
	Photos []struct {
		Src struct {
			Large string `json:"large"`
		} `json:"src"`
	} `json:"photos"`
}

func (a *PexelsAdapterImpl) SearchImage(query string) (string, error) {
	if a.APIKey == "" {
		return "", fmt.Errorf("pexels API key is required")
	}

	params := url.Values{}
	params.Add("query", query)
	params.Add("per_page", "1")
	params.Add("orientation", "square")

	fullURL := fmt.Sprintf("%s/search?%s", a.BaseURL, params.Encode())
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", a.APIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("pexels API returned status: %d", resp.StatusCode)
	}

	var pexResp pexelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&pexResp); err != nil {
		return "", err
	}

	if len(pexResp.Photos) > 0 {
		return pexResp.Photos[0].Src.Large, nil
	}

	return "", fmt.Errorf("no images found for query: %s", query)
}
