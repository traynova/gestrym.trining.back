package usecases

import (
	"fmt"

	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

// GetTrainingPlanUseCase retrieves a single training plan by ID with full nested data.
type GetTrainingPlanUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
}

func NewGetTrainingPlanUseCase(repo interfaces.TrainingPlanRepository) *GetTrainingPlanUseCase {
	return &GetTrainingPlanUseCase{PlanRepo: repo}
}

// Execute fetches the plan and maps it to the frontend-friendly DTO.
// requestingUserID and requestingRoleID are used to enforce access control:
//   - A USER (RoleCliente = 4) may only view plans assigned to themselves.
//   - A TRAINER (RoleCoach = 3) or ADMIN may view any plan.
func (uc *GetTrainingPlanUseCase) Execute(planID uint, requestingUserID uint, requestingRoleID uint) (*dtos.TrainingPlanResponse, error) {
	plan, err := uc.PlanRepo.FindByID(planID)
	if err != nil {
		return nil, fmt.Errorf("error fetching training plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("training plan with ID %d not found", planID)
	}

	// Access control: regular users can only see their own assigned plans
	const roleCliente = uint(4)
	if requestingRoleID == roleCliente {
		if plan.AssignedTo == nil || *plan.AssignedTo != requestingUserID {
			return nil, fmt.Errorf("access denied: this plan is not assigned to you")
		}
	}

	return mapPlanToResponse(plan), nil
}
