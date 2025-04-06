download:
	@go mod download

tidy:
	@go mod tidy

lint:
	@golangci-lint run ./...

tests:
	@go test -cover -race ./...

mockery:
	@mockery --all --case=underscore