#!/bin/sh

go test . -bench BenchmarkSkipList \
    -cpuprofile prof.cpu \
    -memprofile prof.mem \
    -v -benchmem -benchtime 5s

go tool pprof -output cpu.svg  -svg prof.cpu
go tool pprof -output mem.svg  -svg prof.mem
