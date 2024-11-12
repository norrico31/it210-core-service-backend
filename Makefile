build:
	@go build -o bin/it210 cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/it210

run-dev:
	@go run ./cmd

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

GO_CMD=go
SEED_CMD=cmd/seed

build-seeder:
	$(GO_CMD) build -o bin/seeder $(SEED_CMD)/main.go

seed-all: build-seeder
	./bin/seeder

clean:
	rm -rf bin