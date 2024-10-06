package src

import (
	"log"
	"time"
)

type BenchmarkResult struct {
	InitializationDuration time.Duration
	PopulationDuration     time.Duration
	DeletionDuration       time.Duration
	CleanUpDuration        time.Duration
	TotalSizeInBytes       int64
}

func ExecuteBenchmark(b *Benchmark) *BenchmarkResult {
	r := BenchmarkResult{}

	log.Printf("Initializing benchmark")
	r.InitializationDuration = executeTimed(b.Init)

	log.Printf("Populating table")
	r.PopulationDuration = executeTimed(b.Populate)

	r.TotalSizeInBytes = b.GetTotalTableSize()

	log.Printf("Deleting first week data")
	r.DeletionDuration = executeTimed(b.DeleteFirstWeekData)

	log.Printf("Cleaning up")
	r.CleanUpDuration = executeTimed(b.CleanUp)

	return &r
}

func executeTimed(job func()) time.Duration {
	now := time.Now()
	job()
	return time.Since(now)
}

func GetAvgResult(results []*BenchmarkResult) *BenchmarkResult {
	var totalInitDuration, totalPopDuration, totalDelDuration, totalCleanUpDuration time.Duration
	var totalSize int64

	for _, result := range results {
		totalInitDuration += result.InitializationDuration
		totalPopDuration += result.PopulationDuration
		totalDelDuration += result.DeletionDuration
		totalCleanUpDuration += result.CleanUpDuration
		totalSize += result.TotalSizeInBytes
	}

	count := time.Duration(len(results))

	return &BenchmarkResult{
		InitializationDuration: totalInitDuration / count,
		PopulationDuration:     totalPopDuration / count,
		DeletionDuration:       totalDelDuration / count,
		CleanUpDuration:        totalCleanUpDuration / count,
		TotalSizeInBytes:       totalSize / int64(len(results)),
	}
}
