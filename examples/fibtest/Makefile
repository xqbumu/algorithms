test:
	go test -bench=^BenchmarkFibonacci$$ -cpuprofile cpu.out -memprofile mem.out ./...

pprof:
	go tool pprof cpu.out cpu.out

pprof-http:
	go tool pprof -http :8081 cpu.out
