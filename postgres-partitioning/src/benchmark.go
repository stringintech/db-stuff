package src

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Benchmark struct {
	ctx      context.Context
	dbConn   *pgxpool.Pool
	scenario Scenario
}

func NewBenchmark(ctx context.Context, dbConn *pgxpool.Pool, partitionMode ScenarioMode) *Benchmark {
	b := Benchmark{
		ctx:    ctx,
		dbConn: dbConn,
	}

	if partitionMode == Simple {
		b.scenario = &SimpleScenario{
			ctx:    ctx,
			dbConn: dbConn,
		}
	} else {
		b.scenario = &PartitionScenario{
			ctx:               ctx,
			dbConn:            dbConn,
			dropPartitionMode: partitionMode == DropPartition,
		}
	}

	return &b
}

func (b *Benchmark) Init() {
	b.scenario.Init()
}

func (b *Benchmark) Populate() {
	// generating 4 million records; approximately 1 million per week
	_, err := b.dbConn.Exec(b.ctx, `

		insert into records (time, body)
		select
		'2023-01-01 00:00:00'::timestamptz + (n * interval '0.6 seconds') as time,
		repeat('This is a long body of text that is being repeated to increase the length. ', 50) as body
		from generate_series(0, 4000000) as n;

	`)
	if err != nil {
		log.Fatalf("Failed to populate table: %v\n", err)
	}
}

func (b *Benchmark) GetTotalTableSize() int64 {
	return b.scenario.GetTotalSizeInBytes()
}

func (b *Benchmark) DeleteFirstWeekData() {
	b.scenario.DeleteFirstWeekData()
}

func (b *Benchmark) CleanUp() {
	b.scenario.CleanUp()
}
