package main

import (
	"context"
	"log"
	"os"
	"time"

	"gestrym-training/internal/etl/extractor"
	"gestrym-training/internal/etl/loader"
	"gestrym-training/internal/etl/pipeline"
	"gestrym-training/internal/etl/transformer"
	"gestrym-training/src/common/config"
	nutritionAdapters "gestrym-training/src/nutrition/infrastructure/adapters"
	nutritionRepos "gestrym-training/src/nutrition/infrastructure/repositories"
	trainingAdapters "gestrym-training/src/training/infrastructure/adapters"

	"github.com/spf13/viper"
)

func main() {
	log.Printf("[ETL START] Initializing Food ETL Pipeline")

	// 0. Load environment
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "./deployment/env_local.yaml"
	}
	viper.SetConfigFile(configPath)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Fatal error config file: %s", err)
	}

	// 2. Database Connection
	dbConn, err := config.MigrateDB()
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}
	db := dbConn.GetDB()

	// 3. Initialize Adapters
	usdaAdapter := nutritionAdapters.NewUSDAAdapterImpl("", viper.GetString("USDA_API_KEY"))
	pexelsAdapter := nutritionAdapters.NewPexelsAdapterImpl(viper.GetString("PEXELS_API_KEY"))
	legacyStorage := trainingAdapters.NewFileStorageAdapterImpl(viper.GetString("STORAGE_SERVICE_URL"), viper.GetString("STORAGE_SERVICE_API_KEY"))
	storageService := nutritionAdapters.NewStorageServiceAdapterImpl(legacyStorage)
	foodRepo := nutritionRepos.NewFoodRepositoryImpl(db)

	// 4. Initialize ETL Stages
	ext := extractor.NewUSDAExtractor(usdaAdapter)
	trans := transformer.NewFoodTransformer(pexelsAdapter)
	ld := loader.NewDBLoader(foodRepo, storageService)

	// 5. Initialize Pipeline Manager
	manager := pipeline.NewETLManager(ext, trans, ld)

	// 6. Run Execution
	queries := []string{"chicken", "beef", "rice", "egg", "milk", "fish", "potato", "banana", "salmon", "broccoli"}
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	log.Printf("[ETL RUN] Starting worker pool with 5 workers")
	manager.Run(ctx, queries, 5)

	log.Printf("[ETL DONE] Success")
	os.Exit(0)
}
