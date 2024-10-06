package src

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type Scenario interface {
	Init()
	CleanUp()
	GetTotalSizeInBytes() int64
	DeleteFirstWeekData()
}

type SimpleScenario struct {
	ctx    context.Context
	dbConn *pgxpool.Pool
}

func (s *SimpleScenario) Init() {
	_, err := s.dbConn.Exec(s.ctx, `

		create table records (
		    id bigserial primary key,
		    time timestamptz not null,
		    body text
		);

	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}

	_, err = s.dbConn.Exec(s.ctx, `
		create index idx_records_time on records (time);
	`)
	if err != nil {
		log.Fatalf("Failed to create index on table: %v\n", err)
	}
}

func (s *SimpleScenario) CleanUp() {
	_, err := s.dbConn.Exec(s.ctx, `
		drop table records;
	`)
	if err != nil {
		panic(err)
	}
}

func (s *SimpleScenario) GetTotalSizeInBytes() int64 {
	var size int64
	err := s.dbConn.QueryRow(s.ctx, `
		select pg_relation_size('records')
	`).Scan(&size)
	if err != nil {
		log.Fatalf("Failed to get table size: %v\n", err)
	}
	return size
}

func (s *SimpleScenario) DeleteFirstWeekData() {
	_, err := s.dbConn.Exec(s.ctx, `
		delete from records where time < '2023-01-08 00:00:00'
	`)
	if err != nil {
		log.Fatalf("Failed to delete old data: %v\n", err)
	}
}
