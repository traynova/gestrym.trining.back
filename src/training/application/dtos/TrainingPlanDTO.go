package dtos

import "time"

// ─── Request DTOs ────────────────────────────────────────────────────────────

// CreateTrainingPlanRequest is the request body for POST /training-plans.
type CreateTrainingPlanRequest struct {
	Name         string `json:"name"         binding:"required"`
	Description  string `json:"description"`
	DurationDays int    `json:"durationDays" binding:"required,min=1"`
	IsTemplate   bool   `json:"isTemplate"`
}

// AssignTrainingPlanRequest is the request body for POST /training-plans/:id/assign.
type AssignTrainingPlanRequest struct {
	UserID    uint      `json:"userId"    binding:"required"`
	StartDate time.Time `json:"startDate" binding:"required"`
}

// AddTrainingDayRequest is the request body for POST /training-plans/:id/days.
type AddTrainingDayRequest struct {
	DayNumber int    `json:"dayNumber" binding:"required,min=1"`
	WorkoutID uint   `json:"workoutId" binding:"required"`
	Notes     string `json:"notes"`
}

// CloneTrainingPlanRequest is the request body for POST /training-plans/:id/clone.
type CloneTrainingPlanRequest struct {
	TargetUserID uint `json:"targetUserId" binding:"required"`
}

// UpdateDayCompletionRequest is the request body for PATCH /training-plans/:id/days/:dayId/complete.
type UpdateDayCompletionRequest struct {
	IsCompleted bool `json:"isCompleted"`
}

// ─── Response DTOs ───────────────────────────────────────────────────────────

// TrainingPlanResponse is the top-level response with nested days.
type TrainingPlanResponse struct {
	ID           uint                  `json:"id"`
	Name         string                `json:"name"`
	Description  string                `json:"description"`
	DurationDays int                   `json:"durationDays"`
	CreatedBy    uint                  `json:"createdBy"`
	AssignedTo   *uint                 `json:"assignedTo"`
	IsTemplate   bool                  `json:"isTemplate"`
	Days         []TrainingDayResponse `json:"days"`
	CreatedAt    time.Time             `json:"createdAt"`
	UpdatedAt    time.Time             `json:"updatedAt"`
}

// TrainingDayResponse represents a single day in the plan with its workout.
type TrainingDayResponse struct {
	ID          uint                 `json:"id"`
	DayNumber   int                  `json:"dayNumber"`
	Notes       string               `json:"notes"`
	IsCompleted bool                 `json:"isCompleted"`
	Workout     WorkoutSummaryResponse `json:"workout"`
}

// WorkoutSummaryResponse is a lightweight workout view used inside plan day responses.
type WorkoutSummaryResponse struct {
	ID        uint                   `json:"id"`
	Name      string                 `json:"name"`
	Exercises []WorkoutExerciseResponse `json:"exercises"`
}
