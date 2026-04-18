package adapters

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"gestrym-training/src/common/models"
)

type USDAAdapterImpl struct {
	BaseURL string
	APIKey  string
}

func NewUSDAAdapterImpl(baseURL, apiKey string) *USDAAdapterImpl {
	if baseURL == "" {
		baseURL = "https://api.nal.usda.gov/fdc/v1"
	}
	return &USDAAdapterImpl{
		BaseURL: baseURL,
		APIKey:  apiKey,
	}
}

type usdaNutrient struct {
	Name  string  `json:"nutrientName"`
	Value float64 `json:"value"`
}

type usdaFoodItem struct {
	Description   string         `json:"description"`
	FoodCategory  string         `json:"foodCategory"`
	FoodNutrients []usdaNutrient `json:"foodNutrients"`
}

type usdaSearchResponse struct {
	Foods []usdaFoodItem `json:"foods"`
}

func (a *USDAAdapterImpl) SearchFoods(query string) ([]models.Food, error) {
	params := url.Values{}
	params.Add("query", query)
	params.Add("api_key", a.APIKey)
	params.Add("pageSize", "10") // Limit per search to keep it clean

	fullURL := fmt.Sprintf("%s/foods/search?%s", a.BaseURL, params.Encode())
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("USDA API returned status: %d", resp.StatusCode)
	}

	var usdaResp usdaSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&usdaResp); err != nil {
		return nil, err
	}

	var foods []models.Food
	for _, item := range usdaResp.Foods {
		food := models.Food{
			Name: item.Description,
			Category: models.FoodCategory{
				Name: item.FoodCategory,
			},
		}

		// Map Nutrients
		for _, n := range item.FoodNutrients {
			name := strings.ToLower(n.Name)
			if strings.Contains(name, "energy") {
				food.Calories = n.Value
			} else if strings.Contains(name, "protein") {
				food.Protein = n.Value
			} else if strings.Contains(name, "carbohydrate") {
				food.Carbs = n.Value
			} else if strings.Contains(name, "lipid") || strings.Contains(name, "fat") {
				food.Fats = n.Value
			}
		}
		foods = append(foods, food)
	}

	return foods, nil
}
