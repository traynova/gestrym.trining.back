package usecases

import (
	"fmt"

	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

// GetUserTrainingPlansUseCase retrieves all plans assigned to a specific user.
type GetUserTrainingPlansUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
}

func NewGetUserTrainingPlansUseCase(repo interfaces.TrainingPlanRepository) *GetUserTrainingPlansUseCase {
	return &GetUserTrainingPlansUseCase{PlanRepo: repo}
}

// Execute returns all plans assigned to targetUserID.
// A USER (RoleCliente = 4) may only query their own plans.
// A TRAINER (RoleCoach = 3) or ADMIN may query any user's plans.
func (uc *GetUserTrainingPlansUseCase) Execute(targetUserID uint, requestingUserID uint, requestingRoleID uint) ([]dtos.TrainingPlanResponse, error) {
	const roleCliente = uint(4)
	if requestingRoleID == roleCliente && targetUserID != requestingUserID {
		return nil, fmt.Errorf("access denied: you can only view your own training plans")
	}

	plans, err := uc.PlanRepo.FindByUserID(targetUserID)
	if err != nil {
		return nil, fmt.Errorf("could not fetch training plans for user %d: %w", targetUserID, err)
	}

	responses := make([]dtos.TrainingPlanResponse, 0, len(plans))
	for i := range plans {
		responses = append(responses, *mapPlanToResponse(&plans[i]))
	}
	return responses, nil
}
