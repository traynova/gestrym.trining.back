package handlers

import (
	"net/http"
	"strconv"

	"gestrym-training/src/training/application/usecases"
	"github.com/gin-gonic/gin"
)

type WorkoutHandler struct {
	GetFullUseCase *usecases.GetWorkoutFullUseCase
}

func NewWorkoutHandler(getFullUC *usecases.GetWorkoutFullUseCase) *WorkoutHandler {
	return &WorkoutHandler{
		GetFullUseCase: getFullUC,
	}
}

// GetWorkoutFull godoc
// @Summary      Get full workout details
// @Description  Retrieve a workout with its exercises and sets in a frontend-friendly structure.
// @Tags         Workouts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Workout ID"
// @Success      200  {object}  map[string]interface{}
// @Router       /public/workouts/{id}/full [get]
func (h *WorkoutHandler) GetWorkoutFull(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	workout, err := h.GetFullUseCase.Execute(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Workout not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": workout})
}
