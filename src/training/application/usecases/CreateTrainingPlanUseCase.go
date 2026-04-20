package usecases

import (
	"fmt"

	"gestrym-training/src/common/models"
	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

// CreateTrainingPlanUseCase handles the creation of a new training plan.
type CreateTrainingPlanUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
}

func NewCreateTrainingPlanUseCase(repo interfaces.TrainingPlanRepository) *CreateTrainingPlanUseCase {
	return &CreateTrainingPlanUseCase{PlanRepo: repo}
}

// Execute validates the request and persists a new TrainingPlan.
// createdBy is the authenticated user's ID (trainer or regular user).
func (uc *CreateTrainingPlanUseCase) Execute(req dtos.CreateTrainingPlanRequest, createdBy uint) (*dtos.TrainingPlanResponse, error) {
	if req.DurationDays <= 0 {
		return nil, fmt.Errorf("durationDays must be a positive integer")
	}

	plan := &models.TrainingPlan{
		Name:         req.Name,
		Description:  req.Description,
		DurationDays: req.DurationDays,
		CreatedBy:    createdBy,
		IsTemplate:   req.IsTemplate,
	}

	saved, err := uc.PlanRepo.Create(plan)
	if err != nil {
		return nil, fmt.Errorf("could not create training plan: %w", err)
	}

	return mapPlanToResponse(saved), nil
}
