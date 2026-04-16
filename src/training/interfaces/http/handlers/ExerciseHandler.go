package handlers

import (
	"net/http"
	"strconv"

	"gestrym-training/src/training/application/usecases"
	"gestrym-training/src/training/domain/interfaces"
	"github.com/gin-gonic/gin"
)

type ExerciseHandler struct {
	ImportUseCase *usecases.ImportExercisesUseCase
	Repo          interfaces.ExerciseRepository
}

func NewExerciseHandler(importUC *usecases.ImportExercisesUseCase, repo interfaces.ExerciseRepository) *ExerciseHandler {
	return &ExerciseHandler{
		ImportUseCase: importUC,
		Repo:          repo,
	}
}

// ListExercises godoc
// @Summary      Get list of exercises
// @Description  Retrieve all exercises. Optionally filter by bodyPart and target.
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Param        bodyPart  query     string  false  "Filter by body part"
// @Param        target    query     string  false  "Filter by target"
// @Success      200       {object}  map[string]interface{}
// @Failure      500       {object}  map[string]interface{}
// @Router       /public/exercises [get]
// ListExercises handles GET /exercises
func (h *ExerciseHandler) ListExercises(c *gin.Context) {
	bodyPart := c.Query("bodyPart")
	target := c.Query("target")

	exercises, err := h.Repo.ListAll(bodyPart, target)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch exercises"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exercises})
}

// GetExercise godoc
// @Summary      Get exercise by ID
// @Description  Retrieve a specific exercise using its unique ID.
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Exercise ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /public/exercises/{id} [get]
// GetExercise retrieves a specific exercise by ID, handles GET /exercises/:id
func (h *ExerciseHandler) GetExercise(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	exercise, err := h.Repo.FindByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Exercise not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": exercise})
}

// ImportExercises godoc
// @Summary      Manually import exercises
// @Description  Fetch exercises from external API and store them locally. Highly idempotent.
// @Tags         Exercises
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /public/exercises/import [post]
// ImportExercises manually triggers the import via the UseCase, handles POST /exercises/import
func (h *ExerciseHandler) ImportExercises(c *gin.Context) {
	err := h.ImportUseCase.Execute()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to import exercises: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Exercises imported successfully"})
}
