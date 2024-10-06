# PostgreSQL Table Partitioning for Mass Deletions

This project benchmarks the impact of table partitioning on mass deletions in PostgreSQL, comparing performance across different scenarios.

## Quick Start

```bash
git clone https://github.com/stringintech/db-stuff/postgres-partitioning.git
cd postgres-partitioning
go run main.go --mode=[simple|partition|drop-partition]
```

## Benchmark Scenarios

1. **Simple Table:** No partitioning
2. **Partitioned Table (Row Deletion):** Weekly partitions, deleting rows
3. **Partitioned Table (Partition Drop):** Weekly partitions, dropping entire partition

Each scenario:
- Populates 4 million records (1 million per week)
- Measures population and deletion durations
- Records total table size
- Runs 5 iterations for averaged results

## Results Summary

| Scenario | Population | Deletion | Size (Bytes) |
|----------|------------|----------|--------------|
| Simple   | 41.70s     | 1.26s    | 728,178,688  |
| Partition (Delete) | 42.87s | 733.99ms | 908,304,384 |
| Partition (Drop)   | 42.86s | 6.43ms   | 908,304,384 |

## Key Findings

- **Fastest Deletion:** Dropping partitions (6.43ms)
- **Storage Trade-off:** Partitioned tables use ~25% more storage
- **Insertion Impact:** Partitioning slightly increases population time (~2.8%)

## Project Structure

- `benchmark.go`: Core benchmark logic
- `scenario.go`, `partition_scenario.go`: Scenario implementations
- `main.go`: CLI and benchmark execution

## Code Highlights

```go
// Scenario initialization
func NewBenchmark(ctx context.Context, dbConn *pgxpool.Pool, partitionMode ScenarioMode) *Benchmark {
// ... (implementation)
}

// Partitioned deletion
func (s *PartitionScenario) DeleteFirstWeekData() {
q := s.dropPartitionMode ?
`ALTER TABLE records DETACH PARTITION records_partitioned_1; DROP TABLE records_partitioned_1;` :
`DELETE FROM records WHERE time < '2023-01-08 00:00:00'`
// ... (execution)
}
```

## Conclusion

Table partitioning significantly improves mass deletion performance in PostgreSQL, especially when dropping entire partitions. The minor trade-offs in insertion speed and storage usage are often outweighed by the substantial gains in deletion efficiency for large-scale data management scenarios.

For detailed implementation and full benchmark results, explore the repository.