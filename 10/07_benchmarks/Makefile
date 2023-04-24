.PHONY: bench bench-2s bench-1000x benchmem bench-10 bench-cpu install-benchstat benchstats

bench:
	go test -bench=.

bench-2s:
	go test -bench=. -benchtime=2s

bench-1000x:
	go test -bench=. -benchtime=1000x

benchmem:
	go test -bench=. -benchmem

bench-10:
	go test -bench=. -count=10

bench-cpu:
	go test -bench=. -cpu=1,2,4

install-benchstat:
	go install golang.org/x/perf/cmd/benchstat@latest

benchstats:
	go test -bench=. -count=10 | tee stats.txt
	benchstat stats.txt
