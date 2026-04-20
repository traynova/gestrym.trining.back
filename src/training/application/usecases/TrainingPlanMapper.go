package usecases

import (
	"gestrym-training/src/common/models"
	"gestrym-training/src/training/application/dtos"
)

// mapPlanToResponse converts a models.TrainingPlan (with preloaded Days/Workouts) to the frontend DTO.
func mapPlanToResponse(plan *models.TrainingPlan) *dtos.TrainingPlanResponse {
	resp := &dtos.TrainingPlanResponse{
		ID:           plan.ID,
		Name:         plan.Name,
		Description:  plan.Description,
		DurationDays: plan.DurationDays,
		CreatedBy:    plan.CreatedBy,
		AssignedTo:   plan.AssignedTo,
		IsTemplate:   plan.IsTemplate,
		Days:         make([]dtos.TrainingDayResponse, 0, len(plan.Days)),
		CreatedAt:    plan.CreatedAt,
		UpdatedAt:    plan.UpdatedAt,
	}

	for i := range plan.Days {
		resp.Days = append(resp.Days, *mapDayToResponse(&plan.Days[i]))
	}
	return resp
}

// mapDayToResponse converts a models.TrainingDay (with preloaded Workout) to the frontend DTO.
func mapDayToResponse(day *models.TrainingDay) *dtos.TrainingDayResponse {
	resp := &dtos.TrainingDayResponse{
		ID:          day.ID,
		DayNumber:   day.DayNumber,
		Notes:       day.Notes,
		IsCompleted: day.IsCompleted,
		Workout: dtos.WorkoutSummaryResponse{
			ID:        day.Workout.ID,
			Name:      day.Workout.Name,
			Exercises: make([]dtos.WorkoutExerciseResponse, 0, len(day.Workout.Exercises)),
		},
	}

	for _, we := range day.Workout.Exercises {
		resp.Workout.Exercises = append(resp.Workout.Exercises, dtos.WorkoutExerciseResponse{
			ExerciseID: we.ExerciseID,
			Name:       we.Exercise.Name,
			GifURL:     we.Exercise.GifURL,
			Sets:       []dtos.WorkoutSetResponse{}, // populated by GetWorkoutFull when needed
		})
	}
	return resp
}
