package main

import (
	"log"

	"github.com/spf13/viper"
	"gestrym-training/src/common/config"
	"gestrym-training/src/common/models"
	"gestrym-training/src/training/application/usecases"
	"gestrym-training/src/training/infrastructure/adapters"
	"gestrym-training/src/training/infrastructure/repositories"
)

func main() {
	log.Println("Initializing standalone script to import exercises...")

	// 0. Load environment
	viper.SetConfigFile("./deployment/env_local.yaml")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: error reading env file: %v", err)
	}

	config.InitEnvironment(true)
	dbConn := config.NewPostgresConnection()
	db := dbConn.GetDB()

	err := db.AutoMigrate(&models.Exercise{})
	if err != nil {
		log.Fatalf("Failed to auto-migrate Exercise model: %v", err)
	}

	adapter := adapters.NewExerciseDBAdapterImpl("", viper.GetString("RAPID_API_KEY"), viper.GetString("RAPID_API_HOST"))
	repo := repositories.NewExerciseRepositoryImpl(db)
	useCase := usecases.NewImportExercisesUseCase(adapter, repo)

	err = useCase.Execute()
	if err != nil {
		log.Fatalf("Import process failed: %v", err)
	}

	log.Println("Import process completely finished.")
}
