package usecases

import (
	"errors"
	"gestrym-training/src/common/models"
	"gestrym-training/src/training/domain/interfaces"
)

type AdaptTrainingPlanUseCase struct {
	PlanRepo interfaces.TrainingPlanRepository
	DayRepo  interfaces.TrainingDayRepository
}

func NewAdaptTrainingPlanUseCase(
	planRepo interfaces.TrainingPlanRepository,
	dayRepo interfaces.TrainingDayRepository,
) *AdaptTrainingPlanUseCase {
	return &AdaptTrainingPlanUseCase{PlanRepo: planRepo, DayRepo: dayRepo}
}

func (u *AdaptTrainingPlanUseCase) Execute(userID uint) (*models.TrainingPlan, string, error) {
	// 1. Get latest plan
	latestPlan, err := u.PlanRepo.FindLatestByUserID(userID)
	if err != nil {
		return nil, "", err
	}
	if latestPlan == nil {
		return nil, "", errors.New("no plan found for user to adapt")
	}

	// 2. Calculate progress
	totalDays := len(latestPlan.Days)
	if totalDays == 0 {
		return nil, "", errors.New("latest plan has no days to evaluate")
	}

	completedDays := 0
	for _, day := range latestPlan.Days {
		if day.IsCompleted {
			completedDays++
		}
	}

	completionRate := float64(completedDays) / float64(totalDays)

	// 3. Logic for adaptation
	var recommendation string
	var newPlan *models.TrainingPlan

	if completionRate >= 0.8 {
		recommendation = "You're doing great! We've created a more challenging version of your plan."
		// Level up: Clone and mark as adapted
		newPlan = &models.TrainingPlan{
			Name:         latestPlan.Name + " (Adapted - High Intensity)",
			Description:  latestPlan.Description + "\n\nAdaptation: Increased intensity based on excellent progress.",
			DurationDays: latestPlan.DurationDays,
			CreatedBy:    latestPlan.CreatedBy,
			AssignedTo:   &userID,
			IsTemplate:   false,
		}
		
		savedPlan, err := u.PlanRepo.Create(newPlan)
		if err != nil {
			return nil, "", err
		}

		// Clone days
		for _, day := range latestPlan.Days {
			newDay := &models.TrainingDay{
				TrainingPlanID: savedPlan.ID,
				DayNumber:      day.DayNumber,
				WorkoutID:      day.WorkoutID,
				Notes:          day.Notes + " (Focus on progressive overload)",
				IsCompleted:    false,
			}
			_, _ = u.DayRepo.Create(newDay)
		}
		newPlan = savedPlan
	} else if completionRate >= 0.5 {
		recommendation = "Good progress. Keep pushing to complete the remaining days before increasing intensity."
		newPlan = latestPlan
	} else {
		recommendation = "We noticed you're struggling to keep up. Consider taking a rest day or choosing a lighter plan."
		newPlan = latestPlan
	}

	return newPlan, recommendation, nil
}
