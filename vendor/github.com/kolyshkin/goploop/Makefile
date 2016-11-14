all: build

build:
	go build -v

test:
	go test -v .

bench:
	go test -v -run=XXX -bench=. -benchtime=10s

.PHONY: all build test bench
