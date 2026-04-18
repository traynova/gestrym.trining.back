package main

import (
	"log"
	"os"

	"gestrym-training/src/common/config"
	"gestrym-training/src/common/models"
	"gestrym-training/src/training/application/usecases"
	"gestrym-training/src/training/infrastructure/adapters"
	"gestrym-training/src/training/infrastructure/repositories"

	"github.com/spf13/viper"
)

func main() {
	log.Println("Initializing standalone script to import exercises...")

	// 0. Load environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./deployment/env_local.yaml"
	}
	viper.SetConfigFile(configPath)
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
	storageAdapter := adapters.NewFileStorageAdapterImpl(viper.GetString("STORAGE_SERVICE_URL"), viper.GetString("STORAGE_SERVICE_API_KEY"))
	repo := repositories.NewExerciseRepositoryImpl(db)
	useCase := usecases.NewImportExercisesUseCase(adapter, storageAdapter, repo)

	err = useCase.Execute()
	if err != nil {
		log.Fatalf("Import process failed: %v", err)
	}

	log.Println("Import process completely finished.")
}
