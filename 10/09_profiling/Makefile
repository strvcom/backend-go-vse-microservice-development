.PHONY: profiles profile-cpu profile-mem

profiles:
	go test -bench=. -benchmem -cpuprofile cpu.prof -memprofile mem.prof

profile-cpu:
	go tool pprof -http=":8000" ./cpu.prof

profile-mem:
	go tool pprof -http=":8000" ./mem.prof
