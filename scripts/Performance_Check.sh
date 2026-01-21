#!/usr/env/bin bash

echo "This is intended to have a file name passed in as the first argument"
go test -run=^$ -bench=Benchmark_Full_Process | tail -3 > Performance/$1.txt
