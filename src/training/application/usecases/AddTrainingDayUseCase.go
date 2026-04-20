package usecases

import (
	"fmt"

	"gestrym-training/src/common/models"
	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

// AddTrainingDayUseCase adds a new day entry to an existing training plan.
type AddTrainingDayUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
	DayRepo  interfaces.TrainingDayRepository
}

func NewAddTrainingDayUseCase(
	planRepo interfaces.TrainingPlanRepository,
	dayRepo interfaces.TrainingDayRepository,
) *AddTrainingDayUseCase {
	return &AddTrainingDayUseCase{PlanRepo: planRepo, DayRepo: dayRepo}
}

// Execute validates that the plan exists and the day number is within bounds, then persists the entry.
func (uc *AddTrainingDayUseCase) Execute(planID uint, req dtos.AddTrainingDayRequest) (*dtos.TrainingDayResponse, error) {
	plan, err := uc.PlanRepo.FindByID(planID)
	if err != nil {
		return nil, fmt.Errorf("error fetching training plan: %w", err)
	}
	if plan == nil {
		return nil, fmt.Errorf("training plan with ID %d not found", planID)
	}

	// Business rule: day number must be within the plan's duration
	if req.DayNumber < 1 || req.DayNumber > plan.DurationDays {
		return nil, fmt.Errorf("dayNumber %d is out of range for a %d-day plan", req.DayNumber, plan.DurationDays)
	}

	day := &models.TrainingDay{
		TrainingPlanID: planID,
		DayNumber:      req.DayNumber,
		WorkoutID:      req.WorkoutID,
		Notes:          req.Notes,
	}

	saved, err := uc.DayRepo.Create(day)
	if err != nil {
		return nil, fmt.Errorf("could not add training day: %w", err)
	}

	// Reload with preloaded Workout for the response
	savedDay, err := uc.DayRepo.FindByPlanID(planID)
	if err == nil {
		for _, d := range savedDay {
			if d.ID == saved.ID {
				return mapDayToResponse(&d), nil
			}
		}
	}

	// Fallback: return minimal response if preload fails
	return mapDayToResponse(saved), nil
}
