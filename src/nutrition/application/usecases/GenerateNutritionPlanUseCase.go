package usecases

import (
	"errors"
	"gestrym-training/src/common/models"
	"gestrym-training/src/nutrition/domain/interfaces"
)

type GenerateNutritionPlanInput struct {
	UserID        uint    `json:"userId"`
	Weight        float64 `json:"weight"`        // kg
	Height        float64 `json:"height"`        // cm
	Age           int     `json:"age"`
	Gender        string  `json:"gender"`        // male, female
	ActivityLevel float64 `json:"activityLevel"` // 1.2 (sedentary) to 1.9 (heavy)
	Objective     string  `json:"objective"`     // weight_loss, muscle_gain, maintenance
}

type GenerateNutritionPlanUseCase struct {
	repo interfaces.NutritionPlanRepository
}

func NewGenerateNutritionPlanUseCase(repo interfaces.NutritionPlanRepository) *GenerateNutritionPlanUseCase {
	return &GenerateNutritionPlanUseCase{repo: repo}
}

func (u *GenerateNutritionPlanUseCase) Execute(input GenerateNutritionPlanInput) (*models.NutritionPlan, error) {
	if input.UserID == 0 {
		return nil, errors.New("user ID is required")
	}

	// 1. Calculate BMR (Mifflin-St Jeor Equation)
	var bmr float64
	if input.Gender == "male" {
		bmr = (10 * input.Weight) + (6.25 * input.Height) - (5 * float64(input.Age)) + 5
	} else {
		bmr = (10 * input.Weight) + (6.25 * input.Height) - (5 * float64(input.Age)) - 161
	}

	// 2. Calculate TDEE (Total Daily Energy Expenditure)
	tdee := bmr * input.ActivityLevel

	// 3. Adjust based on objective
	var dailyCalories float64
	switch input.Objective {
	case "weight_loss":
		dailyCalories = tdee - 500
	case "muscle_gain":
		dailyCalories = tdee + 300
	default:
		dailyCalories = tdee
	}

	// 4. Calculate Macros
	// Protein: 2g per kg
	proteinGrams := input.Weight * 2.0
	proteinCalories := proteinGrams * 4

	// Fats: 0.8g per kg
	fatsGrams := input.Weight * 0.8
	fatsCalories := fatsGrams * 9

	// Carbs: Rest
	remainingCalories := dailyCalories - proteinCalories - fatsCalories
	if remainingCalories < 0 {
		remainingCalories = 0 // Safety check
	}
	carbsGrams := remainingCalories / 4

	// 5. Create Plan
	plan := &models.NutritionPlan{
		UserID:        input.UserID,
		Name:          "AI Generated Plan - " + input.Objective,
		Objective:     input.Objective,
		DailyCalories: dailyCalories,
		DailyProtein:  proteinGrams,
		DailyCarbs:    carbsGrams,
		DailyFats:     fatsGrams,
		IsActive:      true,
	}

	// 6. Deactivate previous plans and save new one
	_ = u.repo.DeactivateAllForUser(input.UserID)
	err := u.repo.Save(plan)
	if err != nil {
		return nil, err
	}

	return plan, nil
}
