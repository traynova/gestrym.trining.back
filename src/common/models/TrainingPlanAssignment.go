package models

import "time"

// TrainingPlanAssignment tracks which trainer assigned which plan to which user, and when it started.
// Prepared for future AI-generated plan assignment flows.
type TrainingPlanAssignment struct {
	ID             uint         `gorm:"primaryKey"         json:"id"`
	TrainingPlanID uint         `gorm:"not null;index"     json:"trainingPlanId"`
	UserID         uint         `gorm:"not null;index"     json:"userId"`
	AssignedBy     uint         `gorm:"not null;index"     json:"assignedBy"` // Trainer ID
	StartDate      time.Time    `gorm:"not null"           json:"startDate"`
	TrainingPlan   TrainingPlan `gorm:"foreignKey:TrainingPlanID" json:"trainingPlan,omitempty"`
	CreatedAt      time.Time    `json:"createdAt"`
	UpdatedAt      time.Time    `json:"updatedAt"`
}
