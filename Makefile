.PHONY: all build test e2e run migrate clean docker-up docker-down

all: build test

build:
	@echo "==> Building Go Backend Binaries..."
	go build -o bin/api.exe ./cmd/api
	go build -o bin/ingestor.exe ./cmd/ingestor
	go build -o bin/aggregator.exe ./cmd/aggregator
	go build -o bin/migrate.exe ./cmd/migrate
	@echo "==> Building Vue 3 Frontend..."
	cd web && pnpm exec vite build

test:
	@echo "==> Running Go Unit and Contract Tests..."
	go test -v ./...

e2e:
	@echo "==> Running Playwright 11-Viewport E2E Tests..."
	cd e2e && pnpm exec playwright test

migrate:
	@echo "==> Applying Database Migrations..."
	go run ./cmd/migrate -file migrations/0001_initial_schema.up.sql

run: build
	@echo "==> Launching BN RiftStats Global API..."
	./bin/api.exe

docker-up:
	docker-compose up -d --build

docker-down:
	docker-compose down -v

clean:
	rm -rf bin/ web/dist/
