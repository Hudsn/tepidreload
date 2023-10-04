.PHONY: build
build:
	rm -f ./build/*
	go build -o ./build/tepid ./cmd/main.go
