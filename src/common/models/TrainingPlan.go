package models

import "time"

// TrainingPlan represents a structured fitness plan (weekly, monthly, or custom) created by a trainer or user.
type TrainingPlan struct {
	ID           uint          `gorm:"primaryKey"              json:"id"`
	Name         string        `gorm:"size:255;not null"       json:"name"`
	Description  string        `gorm:"type:text"               json:"description"`
	DurationDays int           `gorm:"not null;default:7"      json:"durationDays"` // 7, 30, or custom
	CreatedBy    uint          `gorm:"not null;index"          json:"createdBy"`    // Trainer or User ID
	AssignedTo   *uint         `gorm:"index"                   json:"assignedTo"`   // Nullable → template or self-owned
	IsTemplate   bool          `gorm:"default:false"           json:"isTemplate"`   // Reusable plan base
	Days         []TrainingDay `gorm:"foreignKey:TrainingPlanID" json:"days,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}
