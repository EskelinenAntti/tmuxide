.PHONY: build test install herdr

build:
	go build -o bin/zipln ./zipln

test:
	go test ./...

install:
	go install ./zipln

herdr: build
	herdr plugin link .
