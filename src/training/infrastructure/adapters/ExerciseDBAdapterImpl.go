package adapters

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"gestrym-training/src/common/models"
)

type ExerciseDBAdapterImpl struct {
	BaseURL string
	APIKey  string
	Host    string
}

func NewExerciseDBAdapterImpl(baseURL, apiKey, host string) *ExerciseDBAdapterImpl {
	if baseURL == "" {
		baseURL = "https://exercisedb.p.rapidapi.com"
	}
	return &ExerciseDBAdapterImpl{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Host:    host,
	}
}

// Internal structure matching external API payload closely
type externalExercisePayload struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	BodyPart  string `json:"bodyPart"`
	Target    string `json:"target"`
	Equipment string `json:"equipment"`
	GifUrl    string `json:"gifUrl"`
}

func (a *ExerciseDBAdapterImpl) FetchAllExercises() ([]models.Exercise, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/exercises?limit=3000", a.BaseURL), nil)
	if err != nil {
		return nil, err
	}

	if a.APIKey != "" {
		req.Header.Add("X-RapidAPI-Key", a.APIKey)
	}
	if a.Host != "" {
		req.Header.Add("X-RapidAPI-Host", a.Host)
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch from external API: status %d", res.StatusCode)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var payloads []externalExercisePayload
	if err := json.Unmarshal(body, &payloads); err != nil {
		return nil, err
	}

	var exercises []models.Exercise
	for _, p := range payloads {
		exercises = append(exercises, models.Exercise{
			ExtID:     p.Id,
			Name:      p.Name,
			BodyPart:  p.BodyPart,
			Target:    p.Target,
			Equipment: p.Equipment,
			GifURL:    p.GifUrl,
		})
	}

	return exercises, nil
}
