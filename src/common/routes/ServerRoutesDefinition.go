package routes

import (
	"gestrym-training/docs"
	"gestrym-training/src/common/middleware"
	"gestrym-training/src/common/utils"
	"net/http"
	"sync"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"gestrym-training/src/common/config"
	nutritionUseCases "gestrym-training/src/nutrition/application/usecases"
	nutritionAdapters "gestrym-training/src/nutrition/infrastructure/adapters"
	nutritionRepos "gestrym-training/src/nutrition/infrastructure/repositories"
	nutritionHandlers "gestrym-training/src/nutrition/interfaces/http/handlers"
	"gestrym-training/src/training/application/usecases"
	"gestrym-training/src/training/infrastructure/adapters"
	trainingRepos "gestrym-training/src/training/infrastructure/repositories"
	"gestrym-training/src/training/interfaces/http/handlers"
)

type routesDefinition struct {
	serverGroup    *gin.RouterGroup
	publicGroup    *gin.RouterGroup
	privateGroup   *gin.RouterGroup
	internalGroup  *gin.RouterGroup
	protectedGroup *gin.RouterGroup
	logger         utils.ILogger
}

var (
	routesInstance *routesDefinition
	routesOnce     sync.Once
)

func NewRoutesDefinition(serverInstance *gin.Engine) *routesDefinition {
	routesOnce.Do(func() {
		routesInstance = &routesDefinition{}
		routesInstance.logger = utils.NewLogger()
		docs.SwaggerInfo.Title = "Gestrym Training API"
		docs.SwaggerInfo.Description = "API para el manejo de entrenamientos."
		docs.SwaggerInfo.Version = "1.0"
		docs.SwaggerInfo.BasePath = "/gestrym-training"
		routesInstance.addCORSConfig(serverInstance)
		routesInstance.addRoutes(serverInstance)
	})
	return routesInstance
}

func (r *routesDefinition) addCORSConfig(serverInstance *gin.Engine) {
	corsMiddleware := cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "X-API-Key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	})

	// Aplica el middleware CORS
	serverInstance.Use(corsMiddleware)
}

func (r *routesDefinition) addRoutes(serverInstance *gin.Engine) {
	// Add default routes
	r.addDefaultRoutes(serverInstance)

	// Instantiate DB
	dbConn := config.NewPostgresConnection()
	db := dbConn.GetDB()

	// Repositories
	exerciseRepo := trainingRepos.NewExerciseRepositoryImpl(db)
	workoutRepo := trainingRepos.NewWorkoutRepositoryImpl(db)
	foodRepo := nutritionRepos.NewFoodRepositoryImpl(db)
	trainingPlanRepo := trainingRepos.NewTrainingPlanRepositoryImpl(db)
	trainingDayRepo := trainingRepos.NewTrainingDayRepositoryImpl(db)
	nutritionPlanRepo := nutritionRepos.NewNutritionPlanRepositoryImpl(db)
	assignmentRepo := trainingRepos.NewAssignmentRepositoryImpl(db)

	// Adapters & Services
	exerciseAdapter := adapters.NewExerciseDBAdapterImpl("", viper.GetString("RAPID_API_KEY"), viper.GetString("RAPID_API_HOST"))
	storageAdapter := adapters.NewFileStorageAdapterImpl(viper.GetString("STORAGE_SERVICE_URL"), viper.GetString("STORAGE_SERVICE_API_KEY"))
	usdaAdapter := nutritionAdapters.NewUSDAAdapterImpl("", viper.GetString("USDA_API_KEY"))
	pexelsAdapter := nutritionAdapters.NewPexelsAdapterImpl(viper.GetString("PEXELS_API_KEY"))
	storageService := nutritionAdapters.NewStorageServiceAdapterImpl(storageAdapter)

	// Use Cases
	importExerciseUC := usecases.NewImportExercisesUseCase(exerciseAdapter, storageAdapter, exerciseRepo)
	getWorkoutFullUC := usecases.NewGetWorkoutFullUseCase(workoutRepo)

	searchFoodsUC := nutritionUseCases.NewSearchFoodsUseCase(foodRepo)
	getFoodByIDUC := nutritionUseCases.NewGetFoodByIDUseCase(foodRepo)
	importFoodsUC := nutritionUseCases.NewImportFoodsWithImagesUseCase(foodRepo, usdaAdapter, pexelsAdapter, storageService)

	// Training Plan Use Cases
	createTrainingPlanUC := usecases.NewCreateTrainingPlanUseCase(trainingPlanRepo)
	assignTrainingPlanUC := usecases.NewAssignTrainingPlanUseCase(trainingPlanRepo, trainingDayRepo, assignmentRepo)
	getTrainingPlanUC := usecases.NewGetTrainingPlanUseCase(trainingPlanRepo)
	getUserTrainingPlansUC := usecases.NewGetUserTrainingPlansUseCase(trainingPlanRepo)
	addTrainingDayUC := usecases.NewAddTrainingDayUseCase(trainingPlanRepo, trainingDayRepo)
	cloneTrainingPlanUC := usecases.NewCloneTrainingPlanUseCase(trainingPlanRepo, trainingDayRepo, assignmentRepo)
	updateDayCompletionUC := usecases.NewUpdateDayCompletionUseCase(trainingDayRepo, trainingPlanRepo)
	adaptTrainingPlanUC := usecases.NewAdaptTrainingPlanUseCase(trainingPlanRepo, trainingDayRepo)

	generateNutritionPlanUC := nutritionUseCases.NewGenerateNutritionPlanUseCase(nutritionPlanRepo)

	// Controllers
	exerciseHandler := handlers.NewExerciseHandler(importExerciseUC, exerciseRepo)
	workoutHandler := handlers.NewWorkoutHandler(getWorkoutFullUC)
	foodHandler := nutritionHandlers.NewFoodHandler(searchFoodsUC, getFoodByIDUC, importFoodsUC)
	trainingPlanHandler := handlers.NewTrainingPlanHandler(
		createTrainingPlanUC,
		assignTrainingPlanUC,
		getTrainingPlanUC,
		getUserTrainingPlansUC,
		addTrainingDayUC,
		cloneTrainingPlanUC,
		updateDayCompletionUC,
		adaptTrainingPlanUC,
	)

	nutritionPlanHandler := nutritionHandlers.NewNutritionPlanHandler(generateNutritionPlanUC)

	// Add server group
	r.serverGroup = serverInstance.Group(docs.SwaggerInfo.BasePath)
	r.serverGroup.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Add groups
	r.publicGroup = r.serverGroup.Group("/public")

	// Register Exercise endpoints
	exercisesGroup := r.publicGroup.Group("/exercises")
	{
		exercisesGroup.GET("", exerciseHandler.ListExercises)
		exercisesGroup.GET("/:id", exerciseHandler.GetExercise)
		exercisesGroup.POST("/import", exerciseHandler.ImportExercises)
	}

	// Register Workout endpoints
	workoutsGroup := r.publicGroup.Group("/workouts")
	{
		workoutsGroup.GET("/:id/full", workoutHandler.GetWorkoutFull)
	}

	// Register Nutrition endpoints
	foodsGroup := r.publicGroup.Group("/foods")
	{
		foodsGroup.GET("", foodHandler.SearchFoods)
		foodsGroup.GET("/:id", foodHandler.GetFoodByID)
		foodsGroup.POST("/import", foodHandler.ImportFoods)
	}

	r.privateGroup = r.serverGroup.Group("/private")
	r.protectedGroup = r.serverGroup.Group("/protected")

	// Add middleware to private group
	r.privateGroup.Use(middleware.SetupJWTMiddleware())

	r.protectedGroup.Use(middleware.SetupApiKeyMiddleware())

	// Register Training Plan endpoints (JWT protected, under /private)
	trainingPlansGroup := r.privateGroup.Group("/training-plans")
	{
		trainingPlansGroup.POST("", trainingPlanHandler.CreateTrainingPlan)
		trainingPlansGroup.GET("/:id", trainingPlanHandler.GetTrainingPlan)
		trainingPlansGroup.GET("/user/:userId", trainingPlanHandler.GetUserTrainingPlans)
		trainingPlansGroup.POST("/:id/assign",
			middleware.RequireRoles(middleware.RoleCoach, middleware.RoleAdmin),
			trainingPlanHandler.AssignTrainingPlan,
		)
		trainingPlansGroup.POST("/:id/days", trainingPlanHandler.AddTrainingDay)
		trainingPlansGroup.POST("/:id/clone", trainingPlanHandler.CloneTrainingPlan)
		trainingPlansGroup.PATCH("/:id/days/:dayId/complete", trainingPlanHandler.UpdateDayCompletion)
		trainingPlansGroup.POST("/adapt", trainingPlanHandler.AdaptTrainingPlan)
	}

	// Register Nutrition Plan endpoints (JWT protected)
	nutritionPlansGroup := r.privateGroup.Group("/nutrition-plans")
	{
		nutritionPlansGroup.POST("/generate", nutritionPlanHandler.GenerateNutritionPlan)
	}

	// Add routes to remaining groups
	r.addPublicRoutes()
	r.addInternalRoutes()
	r.addProtectedRoutes()

}

func (r *routesDefinition) addDefaultRoutes(serverInstance *gin.Engine) {

	// Handle root
	serverInstance.GET("/", func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "OK",
			"message": "gestrym-training OK...",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusOK, response)
	})

	// Handle 404
	serverInstance.NoRoute(func(cnx *gin.Context) {
		response := map[string]interface{}{
			"code":    "NOT_FOUND",
			"message": "Resource not found",
			"date":    utils.GetCurrentTime(),
		}

		cnx.JSON(http.StatusNotFound, response)
	})
}

func (r *routesDefinition) addPublicRoutes() {

}

func (r *routesDefinition) addPrivateRoutes() {
	// Additional private routes can be added here when they require
	// handler instances instantiated elsewhere. Currently, training plan
	// routes are registered inline in addRoutes() for scope reasons.
}

func (r *routesDefinition) addInternalRoutes() {

}

func (r *routesDefinition) addProtectedRoutes() {
}
