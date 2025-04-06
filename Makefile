download:
	@go mod download

tidy:
	@go mod tidy

lint:
	@golangci-lint run ./...

fmt:
	@golangci-lint fmt ./...

tests:
	@go test -cover -race ./...

mockery:
	@mockery --all --case=underscore