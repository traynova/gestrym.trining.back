package models

import "time"

type NutritionPlan struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserID      uint      `gorm:"not null;index" json:"userId"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Objective   string    `gorm:"size:100;not null" json:"objective"` // e.g., weight_loss, muscle_gain, maintenance
	DailyCalories float64 `json:"dailyCalories"`
	DailyProtein  float64 `json:"dailyProtein"`
	DailyCarbs    float64 `json:"dailyCarbs"`
	DailyFats     float64 `json:"dailyFats"`
	IsActive    bool      `gorm:"default:true" json:"isActive"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}
