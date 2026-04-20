package usecases

import (
	"fmt"
	"time"

	"gestrym-training/src/common/models"
	"gestrym-training/src/training/domain/interfaces"
)

// AssignTrainingPlanUseCase assigns an existing plan to a user. Only trainers should invoke this.
type AssignTrainingPlanUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
	DayRepo  interfaces.TrainingDayRepository
}

func NewAssignTrainingPlanUseCase(
	planRepo interfaces.TrainingPlanRepository,
	dayRepo interfaces.TrainingDayRepository,
) *AssignTrainingPlanUseCase {
	return &AssignTrainingPlanUseCase{PlanRepo: planRepo, DayRepo: dayRepo}
}

// Execute assigns planID to userID.
// startDate is accepted now and will be persisted in TrainingPlanAssignment in a future iteration.
// trainerID is extracted from the JWT context and validated by the route middleware.
func (uc *AssignTrainingPlanUseCase) Execute(planID uint, userID uint, trainerID uint, _ time.Time) error {
	// Guard: ensure the plan exists before assigning
	plan, err := uc.PlanRepo.FindByID(planID)
	if err != nil {
		return fmt.Errorf("error fetching training plan: %w", err)
	}
	if plan == nil {
		return fmt.Errorf("training plan with ID %d not found", planID)
	}

	// Prevent re-assigning a plan that already belongs to a different user (idempotent for same user)
	if plan.AssignedTo != nil && *plan.AssignedTo != userID {
		// Clone the plan so the original can be reused / templated
		cloned := &models.TrainingPlan{
			Name:         plan.Name,
			Description:  plan.Description,
			DurationDays: plan.DurationDays,
			CreatedBy:    trainerID,
			IsTemplate:   false,
		}
		saved, cloneErr := uc.PlanRepo.Create(cloned)
		if cloneErr != nil {
			return fmt.Errorf("could not clone training plan for new assignment: %w", cloneErr)
		}

		// Copy days to the new plan
		existingDays, dayErr := uc.DayRepo.FindByPlanID(planID)
		if dayErr != nil {
			return fmt.Errorf("could not fetch plan days for cloning: %w", dayErr)
		}
		for _, d := range existingDays {
			newDay := &models.TrainingDay{
				TrainingPlanID: saved.ID,
				DayNumber:      d.DayNumber,
				WorkoutID:      d.WorkoutID,
				Notes:          d.Notes,
			}
			if _, errd := uc.DayRepo.Create(newDay); errd != nil {
				return fmt.Errorf("could not clone training day: %w", errd)
			}
		}

		planID = saved.ID // Assign the cloned plan
	}

	return uc.PlanRepo.AssignToUser(planID, userID)
}
