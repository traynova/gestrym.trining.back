package models

import "time"

type Food struct {
	ID         uint         `gorm:"primaryKey" json:"id"`
	Name       string       `gorm:"size:255;not null;index" json:"name"`
	CategoryID uint         `json:"categoryId"`
	Category   FoodCategory `gorm:"foreignKey:CategoryID" json:"category"`
	Calories     float64      `json:"calories"`
	Protein      float64      `json:"protein"`
	Carbs        float64      `json:"carbs"`
	Fats         float64      `json:"fats"`
	ImageURL     string       `gorm:"size:500" json:"imageUrl"`
	CollectionID string       `gorm:"size:255" json:"collectionId"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
}
