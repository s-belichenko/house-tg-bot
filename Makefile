LINTER_IMAGE:='golangci/golangci-lint:v2.5.0-alpine'
MOCKERY_IMAGE:='vektra/mockery'

download:
	@go mod download

tidy:
	@go mod tidy

lint:
	@echo "Linting..." && docker pull ${LINTER_IMAGE} 1> /dev/null
	@docker run --rm --mount type=bind,src=.,dst=/project --workdir=/project -it ${LINTER_IMAGE} /bin/ash -c "/usr/bin/golangci-lint version; /usr/bin/golangci-lint run ./..."

lint_new_code:
	@echo "Linting..." && docker pull ${LINTER_IMAGE} 1> /dev/null
	@docker run --rm --mount type=bind,src=.,dst=/project --workdir=/project -it ${LINTER_IMAGE} /bin/ash -c "/usr/bin/golangci-lint version; /usr/bin/golangci-lint --new-from-merge-base development run ./..."

fix:
	@echo "Fixing..." && docker pull ${LINTER_IMAGE} 1> /dev/null
	@docker run --rm --mount type=bind,src=.,dst=/project --workdir=/project -it ${LINTER_IMAGE} /bin/ash -c "/usr/bin/golangci-lint version; /usr/bin/golangci-lint run --fix ./..."

fmt:
	@echo "Linting..." && docker pull ${LINTER_IMAGE} 1> /dev/null
	@docker run --rm --mount type=bind,src=.,dst=/project --workdir=/project -it ${LINTER_IMAGE} /bin/ash -c "/usr/bin/golangci-lint version; /usr/bin/golangci-lint fmt ./..."

tests:
	@go test -cover -race ./...

mockery:
	@mockery
#	@echo "Mocking..." && docker pull ${MOCKERY_IMAGE} 1> /dev/null
#	@docker run --rm --mount type=bind,src=.,dst=/project --workdir=/project -it ${MOCKERY_IMAGE} /bin/ash -c "mockery version"
