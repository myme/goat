all: test build

build:
	go build

test:
	go test -v myme.no/goat/src

clean:
	rm -f goat

.PHONY: all build test clean
