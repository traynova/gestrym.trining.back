package pipeline

import (
	"context"
	"log"
	"sync"
	"time"

	"gestrym-training/src/common/models"
	"gestrym-training/internal/etl/extractor"
	"gestrym-training/internal/etl/transformer"
	"gestrym-training/internal/etl/loader"
)

type Job struct {
	Food models.Food
}

type Result struct {
	Food             models.Food
	ExternalImageURL string
	Error            error
}

type ETLManager struct {
	Extractor   *extractor.USDAExtractor
	Transformer *transformer.FoodTransformer
	Loader      *loader.DBLoader
}

func NewETLManager(e *extractor.USDAExtractor, t *transformer.FoodTransformer, l *loader.DBLoader) *ETLManager {
	return &ETLManager{e, t, l}
}

func (m *ETLManager) Run(ctx context.Context, queries []string, workerCount int) {
	jobs := make(chan Job, 100)
	results := make(chan Result, 100)
	var wg sync.WaitGroup
	var submitterWg sync.WaitGroup

	// 1. Start Workers (Transform)
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go m.worker(ctx, &wg, jobs, results)
	}

	// 2. Start Submitting Results to Loader in background
	submitterWg.Add(1)
	go func() {
		defer submitterWg.Done()
		m.resultSubmitter(ctx, results)
	}()

	// 3. Start Extraction
	for _, q := range queries {
		log.Printf("[ETL] Extracting query: %s", q)
		foods, err := m.Extractor.Extract(ctx, q)
		if err != nil {
			log.Printf("[ETL ERROR] Extraction failed for %s: %v", q, err)
			continue
		}

		for _, f := range foods {
			jobs <- Job{Food: f}
		}
	}

	// 4. Cleanup
	close(jobs)
	wg.Wait()      // Wait for all workers to finish
	close(results) // Signal submitter to stop
	submitterWg.Wait() // Wait for submitter to finish
	
	log.Printf("[ETL SUCCESS] Pipeline finished")
}

func (m *ETLManager) worker(ctx context.Context, wg *sync.WaitGroup, jobs <-chan Job, results chan<- Result) {
	defer wg.Done()
	for job := range jobs {
		// Retry logic: 3 attempts
		var err error
		var transformedFood models.Food
		var imageURL string

		for attempt := 1; attempt <= 3; attempt++ {
			transformedFood, imageURL, err = m.Transformer.Transform(ctx, job.Food)
			if err == nil {
				break
			}
			log.Printf("[ETL RETRY] Attempt %d failed for %s: %v", attempt, job.Food.Name, err)
			time.Sleep(time.Duration(attempt) * time.Second)
		}

		results <- Result{
			Food:             transformedFood,
			ExternalImageURL: imageURL,
			Error:            err,
		}
	}
}

func (m *ETLManager) resultSubmitter(ctx context.Context, results <-chan Result) {
	for res := range results {
		if res.Error != nil {
			log.Printf("[ETL ERROR] Skipping %s due to permanent failure: %v", res.Food.Name, res.Error)
			continue
		}

		err := m.Loader.Load(ctx, res.Food, res.ExternalImageURL)
		if err != nil {
			log.Printf("[ETL ERROR] Failed to load %s into DB: %v", res.Food.Name, err)
		} else {
			log.Printf("[ETL OK] Loaded: %s", res.Food.Name)
		}
	}
}
