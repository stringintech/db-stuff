package src

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type PartitionScenario struct {
	ctx               context.Context
	dbConn            *pgxpool.Pool
	dropPartitionMode bool
}

func (s *PartitionScenario) Init() {
	_, err := s.dbConn.Exec(s.ctx, `
		create table records (
			id bigserial,
			time timestamptz not null,
			body text
		) partition by range (time);
	`)
	if err != nil {
		log.Fatalf("Failed to create partitioned table: %v\n", err)
	}

	// Create partitions for weekly intervals
	_, err = s.dbConn.Exec(s.ctx, `
		create table records_partitioned_1 partition of records
		for values from ('2023-01-01 00:00:00') to ('2023-01-08 00:00:00');

		create table records_partitioned_2 partition of records
		for values from ('2023-01-08 00:00:00') to ('2023-01-15 00:00:00');

		create table records_partitioned_3 partition of records
		for values from ('2023-01-15 00:00:00') to ('2023-01-22 00:00:00');

		create table records_partitioned_4 partition of records
		for values from ('2023-01-22 00:00:00') to ('2023-02-01 00:00:00');
	`)
	if err != nil {
		log.Fatalf("Failed to create partitions: %v\n", err)
	}

	_, err = s.dbConn.Exec(s.ctx, `
		alter table records_partitioned_1 add primary key (id);
		alter table records_partitioned_2 add primary key (id);
		alter table records_partitioned_3 add primary key (id);
		alter table records_partitioned_4 add primary key (id);
	`)
	if err != nil {
		log.Fatalf("Failed to set primary key on partitions: %v\n", err)
	}

	_, err = s.dbConn.Exec(s.ctx, `
		create index idx_records_partitioned_time_1 on records_partitioned_1 (time);
		create index idx_records_partitioned_time_2 on records_partitioned_2 (time);
		create index idx_records_partitioned_time_3 on records_partitioned_3 (time);
		create index idx_records_partitioned_time_4 on records_partitioned_4 (time);
	`)
	if err != nil {
		log.Fatalf("Failed to create indexes on partitions: %v\n", err)
	}
}

func (s *PartitionScenario) CleanUp() {
	_, err := s.dbConn.Exec(s.ctx, `
		drop table records cascade;
	`)
	if err != nil {
		log.Fatalf("Failed to drop partitioned table: %v\n", err)
	}
}

func (s *PartitionScenario) GetTotalSizeInBytes() int64 {
	var size int64
	err := s.dbConn.QueryRow(s.ctx, `
		select pg_total_relation_size('records_partitioned_1')+
               pg_total_relation_size('records_partitioned_2')+
               pg_total_relation_size('records_partitioned_3')+
               pg_total_relation_size('records_partitioned_4');
	`).Scan(&size)
	if err != nil {
		log.Fatalf("Failed to get partitioned table size: %v\n", err)
	}
	return size
}

func (s *PartitionScenario) DeleteFirstWeekData() {
	q := `delete from records where time < '2023-01-08 00:00:00'`
	if s.dropPartitionMode {
		q = `alter table records detach partition records_partitioned_1;
				drop table records_partitioned_1;`
	}
	_, err := s.dbConn.Exec(s.ctx, q)
	if err != nil {
		log.Fatalf("Failed to delete old data: %v\n", err)
	}
}
