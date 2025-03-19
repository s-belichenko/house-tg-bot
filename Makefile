download:
	@go mod download

tidy:
	@go mod tidy

tests:
	@go test -cover -race ./...

mockery:
	@mockery --all --case=underscore