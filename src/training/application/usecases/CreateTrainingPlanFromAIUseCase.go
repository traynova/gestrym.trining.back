package usecases

import (
	"fmt"
	"gestrym-training/src/common/models"
	"gestrym-training/src/common/utils"
	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

// CreateTrainingPlanFromAIUseCase orchestrates the storage of an AI-generated training plan.
type CreateTrainingPlanFromAIUseCase struct {
	PlanRepo       interfaces.TrainingPlanRepository
	DayRepo        interfaces.TrainingDayRepository
	WorkoutRepo    interfaces.WorkoutRepository
	AssignmentRepo interfaces.AssignmentRepository
	ExerciseRepo   interfaces.ExerciseRepository
}

func NewCreateTrainingPlanFromAIUseCase(
	planRepo interfaces.TrainingPlanRepository,
	dayRepo interfaces.TrainingDayRepository,
	workoutRepo interfaces.WorkoutRepository,
	assignmentRepo interfaces.AssignmentRepository,
	exerciseRepo interfaces.ExerciseRepository,
) *CreateTrainingPlanFromAIUseCase {
	return &CreateTrainingPlanFromAIUseCase{
		PlanRepo:       planRepo,
		DayRepo:        dayRepo,
		WorkoutRepo:    workoutRepo,
		AssignmentRepo: assignmentRepo,
		ExerciseRepo:   exerciseRepo,
	}
}

// Execute validates the input and creates the plan, workouts, days and assignment.
func (uc *CreateTrainingPlanFromAIUseCase) Execute(req dtos.CreateTrainingPlanFromAIRequest) (uint, error) {
	// 1. Validation Logic
	// Validate exercises exist and dayNumber is valid
	exerciseMap := make(map[uint]bool)
	for _, day := range req.Days {
		if day.DayNumber < 1 || day.DayNumber > req.DurationDays {
			return 0, fmt.Errorf("invalid dayNumber %d for plan duration %d", day.DayNumber, req.DurationDays)
		}
		for _, ex := range day.Workout.Exercises {
			if _, ok := exerciseMap[ex.ExerciseID]; !ok {
				exercise, err := uc.ExerciseRepo.FindByID(ex.ExerciseID)
				if err != nil || exercise == nil {
					return 0, fmt.Errorf("exercise with ID %d does not exist", ex.ExerciseID)
				}
				exerciseMap[ex.ExerciseID] = true
			}
		}
	}

	// 2. Create TrainingPlan
	plan := &models.TrainingPlan{
		Name:          req.Name,
		DurationDays:  req.DurationDays,
		AssignedTo:    &req.UserID,
		IsAIGenerated: true,
		CreatedBy:     0, // System/AI generated
	}

	savedPlan, err := uc.PlanRepo.Create(plan)
	if err != nil {
		return 0, fmt.Errorf("failed to create training plan: %w", err)
	}

	// 3 & 4. Create Workouts and TrainingDays
	for _, dayReq := range req.Days {
		workout := &models.Workout{
			UserID: req.UserID,
			Name:   dayReq.Workout.Name,
		}

		for i, exReq := range dayReq.Workout.Exercises {
			workoutEx := models.WorkoutExercise{
				ExerciseID: exReq.ExerciseID,
				Order:      i,
			}
			for _, setReq := range exReq.Sets {
				workoutEx.Sets = append(workoutEx.Sets, models.WorkoutSet{
					Reps:        setReq.Reps,
					RestSeconds: setReq.Rest,
					Type:        "normal",
				})
			}
			workout.Exercises = append(workout.Exercises, workoutEx)
		}

		if err := uc.WorkoutRepo.Create(workout); err != nil {
			return 0, fmt.Errorf("failed to create workout for day %d: %w", dayReq.DayNumber, err)
		}

		trainingDay := &models.TrainingDay{
			TrainingPlanID: savedPlan.ID,
			DayNumber:      dayReq.DayNumber,
			WorkoutID:      workout.ID,
		}

		if _, err := uc.DayRepo.Create(trainingDay); err != nil {
			return 0, fmt.Errorf("failed to create training day %d: %w", dayReq.DayNumber, err)
		}
	}

	// 5. Assign plan to user (history tracking)
	assignment := &models.TrainingPlanAssignment{
		UserID:         req.UserID,
		TrainingPlanID: savedPlan.ID,
		AssignedBy:     0, // AI
		StartDate:      utils.GetCurrentTime(),
	}

	if _, err := uc.AssignmentRepo.Assign(assignment); err != nil {
		return 0, fmt.Errorf("failed to assign plan to user: %w", err)
	}

	return savedPlan.ID, nil
}
