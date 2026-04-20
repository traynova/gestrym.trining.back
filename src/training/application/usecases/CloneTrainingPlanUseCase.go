package usecases

import (
	"fmt"

	"gestrym-training/src/common/models"
	"gestrym-training/src/training/domain/interfaces"
)

// CloneTrainingPlanUseCase handles cloning a plan (template or existing) to a specific user.
type CloneTrainingPlanUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
	DayRepo  interfaces.TrainingDayRepository
}

func NewCloneTrainingPlanUseCase(
	planRepo interfaces.TrainingPlanRepository,
	dayRepo interfaces.TrainingDayRepository,
) *CloneTrainingPlanUseCase {
	return &CloneTrainingPlanUseCase{PlanRepo: planRepo, DayRepo: dayRepo}
}

// Execute performs a deep copy of a training plan and assigns it to targetUserID.
func (uc *CloneTrainingPlanUseCase) Execute(planID uint, targetUserID uint, creatorID uint) (uint, error) {
	// 1. Fetch source plan with its days
	sourcePlan, err := uc.PlanRepo.FindByID(planID)
	if err != nil {
		return 0, fmt.Errorf("error fetching source plan: %w", err)
	}
	if sourcePlan == nil {
		return 0, fmt.Errorf("source plan with ID %d not found", planID)
	}

	// 2. Create the new plan instance
	newPlan := &models.TrainingPlan{
		Name:         fmt.Sprintf("%s (Clone)", sourcePlan.Name),
		Description:  sourcePlan.Description,
		DurationDays: sourcePlan.DurationDays,
		CreatedBy:    creatorID,
		AssignedTo:   &targetUserID,
		IsTemplate:   false,
	}

	savedPlan, err := uc.PlanRepo.Create(newPlan)
	if err != nil {
		return 0, fmt.Errorf("could not create cloned plan: %w", err)
	}

	// 3. Clone all days
	for _, day := range sourcePlan.Days {
		newDay := &models.TrainingDay{
			TrainingPlanID: savedPlan.ID,
			DayNumber:      day.DayNumber,
			WorkoutID:      day.WorkoutID,
			Notes:          day.Notes,
			IsCompleted:    false, // Clones start as uncompleted
		}
		if _, err := uc.DayRepo.Create(newDay); err != nil {
			return 0, fmt.Errorf("could not clone training day %d: %w", day.DayNumber, err)
		}
	}

	return savedPlan.ID, nil
}
