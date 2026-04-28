package handlers

import (
	"net/http"
	"strconv"
	"time"

	"gestrym-training/src/training/application/dtos"
	"gestrym-training/src/training/application/usecases"
	"github.com/gin-gonic/gin"
)

// TrainingPlanHandler wires all training plan endpoints to their use cases.
type TrainingPlanHandler struct {
	CreateUC              *usecases.CreateTrainingPlanUseCase
	AssignUC              *usecases.AssignTrainingPlanUseCase
	GetUC                 *usecases.GetTrainingPlanUseCase
	GetUserUC             *usecases.GetUserTrainingPlansUseCase
	AddDayUC              *usecases.AddTrainingDayUseCase
	CloneUC               *usecases.CloneTrainingPlanUseCase
	UpdateDayCompletionUC *usecases.UpdateDayCompletionUseCase
	AdaptPlanUC           *usecases.AdaptTrainingPlanUseCase
	CreateFromAIUC        *usecases.CreateTrainingPlanFromAIUseCase
}

func NewTrainingPlanHandler(
	createUC *usecases.CreateTrainingPlanUseCase,
	assignUC *usecases.AssignTrainingPlanUseCase,
	getUC *usecases.GetTrainingPlanUseCase,
	getUserUC *usecases.GetUserTrainingPlansUseCase,
	addDayUC *usecases.AddTrainingDayUseCase,
	cloneUC *usecases.CloneTrainingPlanUseCase,
	updateDayCompletionUC *usecases.UpdateDayCompletionUseCase,
	adaptPlanUC *usecases.AdaptTrainingPlanUseCase,
	createFromAIUC *usecases.CreateTrainingPlanFromAIUseCase,
) *TrainingPlanHandler {
	return &TrainingPlanHandler{
		CreateUC:              createUC,
		AssignUC:              assignUC,
		GetUC:                 getUC,
		GetUserUC:             getUserUC,
		AddDayUC:              addDayUC,
		CloneUC:               cloneUC,
		UpdateDayCompletionUC: updateDayCompletionUC,
		AdaptPlanUC:           adaptPlanUC,
		CreateFromAIUC:        createFromAIUC,
	}
}

// ─── Helpers ─────────────────────────────────────────────────────────────────

func getUserIDFromCtx(c *gin.Context) (uint, bool) {
	v, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

func getRoleIDFromCtx(c *gin.Context) (uint, bool) {
	v, exists := c.Get("role_id")
	if !exists {
		return 0, false
	}
	id, ok := v.(uint)
	return id, ok
}

// ─── Handlers ────────────────────────────────────────────────────────────────

// CreateTrainingPlan godoc
// @Summary      Create a training plan
// @Description  Creates a new training plan. Accessible by trainers (role 3) and regular users (role 4).
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      dtos.CreateTrainingPlanRequest  true  "Plan data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      401      {object}  map[string]interface{}
// @Router       /private/training-plans [post]
func (h *TrainingPlanHandler) CreateTrainingPlan(c *gin.Context) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	var req dtos.CreateTrainingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := h.CreateUC.Execute(req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": plan})
}

// GetTrainingPlan godoc
// @Summary      Get training plan by ID
// @Description  Retrieves a training plan with its days and workouts. Users can only see their assigned plans.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Training Plan ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /private/training-plans/{id} [get]
func (h *TrainingPlanHandler) GetTrainingPlan(c *gin.Context) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	roleID, _ := getRoleIDFromCtx(c)

	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	plan, err := h.GetUC.Execute(uint(id), userID, roleID)
	if err != nil {
		if err.Error() == "access denied: this plan is not assigned to you" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plan})
}

// GetUserTrainingPlans godoc
// @Summary      Get training plans for a user
// @Description  Returns all training plans assigned to a specific user. Regular users can only query their own.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        userId  path      int  true  "Target User ID"
// @Success      200     {object}  map[string]interface{}
// @Failure      400     {object}  map[string]interface{}
// @Failure      403     {object}  map[string]interface{}
// @Router       /private/training-plans/user/{userId} [get]
func (h *TrainingPlanHandler) GetUserTrainingPlans(c *gin.Context) {
	requestingUserID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	roleID, _ := getRoleIDFromCtx(c)

	targetID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid userId"})
		return
	}

	plans, err := h.GetUserUC.Execute(uint(targetID), requestingUserID, roleID)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": plans})
}

// AssignTrainingPlan godoc
// @Summary      Assign training plan to a user (TRAINER ONLY)
// @Description  Assigns an existing training plan to a user. Only accessible by trainers (role 3).
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                               true  "Training Plan ID"
// @Param        request  body      dtos.AssignTrainingPlanRequest    true  "Assignment data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Failure      403      {object}  map[string]interface{}
// @Router       /private/training-plans/{id}/assign [post]
func (h *TrainingPlanHandler) AssignTrainingPlan(c *gin.Context) {
	trainerID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var req dtos.AssignTrainingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startDate := req.StartDate
	if startDate.IsZero() {
		startDate = time.Now()
	}

	if err := h.AssignUC.Execute(uint(planID), req.UserID, trainerID, startDate); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "training plan successfully assigned"})
}

// AddTrainingDay godoc
// @Summary      Add a day to a training plan
// @Description  Adds a workout day to an existing training plan. DayNumber must be within the plan's duration.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "Training Plan ID"
// @Param        request  body      dtos.AddTrainingDayRequest   true  "Day data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /private/training-plans/{id}/days [post]
func (h *TrainingPlanHandler) AddTrainingDay(c *gin.Context) {
	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var req dtos.AddTrainingDayRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	day, err := h.AddDayUC.Execute(uint(planID), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": day})
}

// CloneTrainingPlan godoc
// @Summary      Clone a training plan for a user
// @Description  Clones an existing plan (template or other) and assigns it to a specific user.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                          true  "Plan ID to clone"
// @Param        request  body      dtos.CloneTrainingPlanRequest true  "Clone data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /private/training-plans/:id/clone [post]
func (h *TrainingPlanHandler) CloneTrainingPlan(c *gin.Context) {
	creatorID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	var req dtos.CloneTrainingPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newID, err := h.CloneUC.Execute(uint(planID), req.TargetUserID, creatorID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": newID}})
}

// UpdateDayCompletion godoc
// @Summary      Update training day progress
// @Description  Marks a training day as completed or not.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      int                             true  "Plan ID"
// @Param        dayId    path      int                             true  "Day ID"
// @Param        request  body      dtos.UpdateDayCompletionRequest true  "Progress data"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /private/training-plans/:id/days/:dayId/complete [patch]
func (h *TrainingPlanHandler) UpdateDayCompletion(c *gin.Context) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}
	roleID, _ := getRoleIDFromCtx(c)

	planID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid plan ID"})
		return
	}

	dayID, err := strconv.ParseUint(c.Param("dayId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid day ID"})
		return
	}

	var req dtos.UpdateDayCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.UpdateDayCompletionUC.Execute(uint(dayID), uint(planID), userID, roleID, req.IsCompleted); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "progress updated successfully"})
}

// AdaptTrainingPlan godoc
// @Summary      Adapt training plan based on progress
// @Description  Analyzes user progress and creates an adapted version of the latest plan if completion is high.
// @Tags         TrainingPlans
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Router       /private/training-plans/adapt [post]
func (h *TrainingPlanHandler) AdaptTrainingPlan(c *gin.Context) {
	userID, ok := getUserIDFromCtx(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		return
	}

	plan, recommendation, err := h.AdaptPlanUC.Execute(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":           plan,
		"recommendation": recommendation,
	})
}

// CreateFromAI godoc
// @Summary      Create a training plan from AI (INTERNAL ONLY)
// @Description  Stores a plan generated by the AI service. Validates exercises and assigns to user.
// @Tags         Internal
// @Accept       json
// @Produce      json
// @Param        request  body      dtos.CreateTrainingPlanFromAIRequest  true  "AI Plan data"
// @Success      201      {object}  map[string]interface{}
// @Failure      400      {object}  map[string]interface{}
// @Router       /internal/training-plans/ai [post]
func (h *TrainingPlanHandler) CreateFromAI(c *gin.Context) {
	var req dtos.CreateTrainingPlanFromAIRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	planID, err := h.CreateFromAIUC.Execute(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": gin.H{"id": planID}})
}
