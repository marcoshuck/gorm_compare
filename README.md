# GORM Compare
This repository has multiple benchmarks to test gorm's speed using the different strategies provided in their package. All the tests were configured to run the same process with 10 different entities.

## Create benchmark
`creates_test.go` has different benchmarks to get statistics about the creation of a certain entity using single queries vs transactions.

### Results
| Name | Runs | Result |
| ---- | ---- | ------ |
| BenchmarkSingleQueryCreate | 765 | 1546409 ns/op |
| BenchmarkMultiQueryCreate | 793 | 1539416 ns/op |
| BenchmarkMultiQueryCreateTx | 1173 | 966259 ns/op |

## Update benchmark
`update_test.go` has two benchmarks to test the difference between updating an entity without previously checking that the entity exists, and updating an entity that was previously gotten from the database.

### Results
| Name | Runs | Result |
| ---- | ---- | ------ |
| BenchmarkMultiQueryUpdate | 103 | 11770946 ns/op |
| BenchmarkSingleQueryUpdate | 159 | 7438884 ns/op |