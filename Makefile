.PHONY: generate
generate:
	@go generate ./...

.PHONY: build
build: 
	@go build -o bin/ -ldflags "-s -w" -trimpath ./...
