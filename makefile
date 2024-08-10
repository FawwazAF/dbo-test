pkgs = $(shell go list ./... | grep -v vendor | grep -v mocks)

.PHONY: build

build:
	@echo "Reformatting..."
	@go fmt ./...
	@echo "BUILD PROJECT..."
	@go build -v -o ./bin/dbo-test ./cmd/main.go

run-http:
	@echo "RUN HTTP..."
	make build
	@./bin/dbo-test

test:
	@echo "RUN TESTING..."
	@go clean -testcache
	@go test -v -p=2 -cover -race $(pkgs) -coverprofile coverage.out

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down