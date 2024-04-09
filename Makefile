all: test build

build:
	go build

test:
	go test -v github.com/myme/goat/src

clean:
	rm -f goat

.PHONY: all build test clean
