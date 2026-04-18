package usecases

import (
	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/domain/interfaces"
)

type GetWorkoutFullUseCase struct {
	Repo interfaces.WorkoutRepository
}

func NewGetWorkoutFullUseCase(repo interfaces.WorkoutRepository) *GetWorkoutFullUseCase {
	return &GetWorkoutFullUseCase{Repo: repo}
}

func (uc *GetWorkoutFullUseCase) Execute(workoutID uint) (*dtos.WorkoutFullResponse, error) {
	workout, err := uc.Repo.FindFullWorkoutByID(workoutID)
	if err != nil {
		return nil, err
	}

	// Map to DTO
	response := &dtos.WorkoutFullResponse{
		WorkoutID: workout.ID,
		Name:      workout.Name,
		Exercises: make([]dtos.WorkoutExerciseResponse, 0),
	}

	for _, we := range workout.Exercises {
		exerciseDTO := dtos.WorkoutExerciseResponse{
			ExerciseID: we.ExerciseID,
			Name:       we.Exercise.Name,
			GifURL:     we.Exercise.GifURL,
			Sets:       make([]dtos.WorkoutSetResponse, 0),
		}

		for _, s := range we.Sets {
			exerciseDTO.Sets = append(exerciseDTO.Sets, dtos.WorkoutSetResponse{
				Type:        s.Type,
				Reps:        s.Reps,
				Weight:      s.Weight,
				RestSeconds: s.RestSeconds,
			})
		}

		response.Exercises = append(response.Exercises, exerciseDTO)
	}

	return response, nil
}
