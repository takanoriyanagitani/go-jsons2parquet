#!/bin/sh

tinygo \
	build \
	-o ./logs2parquet.wasm \
	-target=wasip1 \
	-opt=z \
	-no-debug \
	./main.go
