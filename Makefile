all: build

BIN = docker-volume-ploop

build:
	go build -v

test:
	go test -v .

clean:
	rm -f $(BIN)

.PHONY: all build test clean
