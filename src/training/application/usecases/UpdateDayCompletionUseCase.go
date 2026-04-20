package usecases

import (
	"fmt"

	"gestrym-training/src/training/domain/interfaces"
)

// UpdateDayCompletionUseCase handles marking a training day as completed or not.
type UpdateDayCompletionUseCase struct {
	DayRepo  interfaces.TrainingDayRepository
	PlanRepo interfaces.TrainingPlanRepository
}

func NewUpdateDayCompletionUseCase(
	dayRepo interfaces.TrainingDayRepository,
	planRepo interfaces.TrainingPlanRepository,
) *UpdateDayCompletionUseCase {
	return &UpdateDayCompletionUseCase{DayRepo: dayRepo, PlanRepo: planRepo}
}

// Execute updates the completion status of a specific day after verifying access.
func (uc *UpdateDayCompletionUseCase) Execute(dayID uint, planID uint, userID uint, roleID uint, isCompleted bool) error {
	// 1. Verify day belongs to the plan
	day, err := uc.DayRepo.FindByID(dayID)
	if err != nil {
		return fmt.Errorf("error fetching training day: %w", err)
	}
	if day == nil || day.TrainingPlanID != planID {
		return fmt.Errorf("training day not found in this plan")
	}

	// 2. Access control: only the assigned user or a trainer/admin can update progress
	plan, err := uc.PlanRepo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("error verifying plan ownership: %w", err)
	}
	if plan == nil {
		return fmt.Errorf("plan not found")
	}

	const roleCliente = uint(4)
	if roleID == roleCliente && (plan.AssignedTo == nil || *plan.AssignedTo != userID) {
		return fmt.Errorf("access denied: you can only update your own plan progress")
	}

	// 3. Update status
	return uc.DayRepo.UpdateCompletionStatus(dayID, isCompleted)
}
