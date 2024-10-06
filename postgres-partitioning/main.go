package main

import (
	"context"
	"flag"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stringintech/db-stuff/postgres-partitioning/src"
	"log"
)

func main() {
	var mode src.ScenarioMode
	flag.Var(&mode, "mode", "set mode (simple, partition, drop-partition)")
	flag.Parse()

	ctx := context.Background()
	connStr := "postgres://postgres:postgres@localhost:6432/benchmark_db"

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer pool.Close()

	log.Printf("Starting benchmark in %s mode", mode.String())

	var results []*src.BenchmarkResult
	for i := 1; i <= 1; i++ {
		log.Printf("Running iteration %d", i)
		b := src.NewBenchmark(ctx, pool, mode)
		result := src.ExecuteBenchmark(b)
		results = append(results, result)
	}

	avg := src.GetAvgResult(results)
	log.Printf("Average Initialization Duration: %s", avg.InitializationDuration)
	log.Printf("Average Population Duration: %s", avg.PopulationDuration)
	log.Printf("Average Deletion Duration: %s", avg.DeletionDuration)
	log.Printf("Average Clean-Up Duration: %s", avg.CleanUpDuration)
	log.Printf("Average Total Size in Bytes: %d", avg.TotalSizeInBytes)
}
