.PHONY: build test install herdr

build:
	go build -o bin/ide ./ide

test:
	go test ./...

install:
	go install ./ide

herdr: build
	herdr plugin link .
