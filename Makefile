.PHONY: deps

deps:
	go mod download && go mod tidy
