tests:
	@go test -cover -race ./...

tidy:
	@go mod tidy