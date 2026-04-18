package dtos

type WorkoutFullResponse struct {
	WorkoutID uint                   `json:"workoutId"`
	Name      string                 `json:"name"`
	Exercises []WorkoutExerciseResponse `json:"exercises"`
}

type WorkoutExerciseResponse struct {
	ExerciseID uint                 `json:"exerciseId"`
	Name       string               `json:"name"`
	GifURL     string               `json:"gifUrl"`
	Sets       []WorkoutSetResponse `json:"sets"`
}

type WorkoutSetResponse struct {
	Type        string  `json:"type"`
	Reps        int     `json:"reps"`
	Weight      float64 `json:"weight"`
	RestSeconds int     `json:"restSeconds"`
}
