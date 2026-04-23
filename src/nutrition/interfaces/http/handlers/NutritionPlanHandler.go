package handlers

import (
	"gestrym-training/src/nutrition/application/usecases"
	"net/http"

	"github.com/gin-gonic/gin"
)

type NutritionPlanHandler struct {
	generateUC *usecases.GenerateNutritionPlanUseCase
}

func NewNutritionPlanHandler(generateUC *usecases.GenerateNutritionPlanUseCase) *NutritionPlanHandler {
	return &NutritionPlanHandler{generateUC: generateUC}
}

// GenerateNutritionPlan handles the POST request to generate a new nutritional plan.
// @Summary Generate nutritional plan
// @Description Calculates macros and creates a plan based on user metrics and goals.
// @Tags nutrition
// @Accept json
// @Produce json
// @Param body body usecases.GenerateNutritionPlanInput true "Plan generation data"
// @Success 200 {object} models.NutritionPlan
// @Router /private/nutrition-plans/generate [post]
func (h *NutritionPlanHandler) GenerateNutritionPlan(c *gin.Context) {
	var input usecases.GenerateNutritionPlanInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// In a real scenario, we'd get UserID from JWT.
	// For now, if not provided in JSON, we can try to extract it from context if middleware is used.
	// if input.UserID == 0 { ... }

	plan, err := h.generateUC.Execute(input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plan)
}
