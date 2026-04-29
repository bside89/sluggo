.DEFAULT_GOAL := help

.PHONY: help run-local run-docker docs build

help:
	@echo "Available commands:"
	@echo "  make run-local       Start PostgreSQL and run the API locally"
	@echo "  make run-docker      Build and start all services with Docker Compose"
	@echo "  make docs            Regenerate Swagger documentation"
	@echo "  make build           Build the API binary"

run-local:
	@bash ./scripts/run-local.sh

run-docker:
	@bash ./scripts/run-docker.sh

docs:
	@bash ./scripts/update-swagger.sh

build:
	@go build -o sluggo cmd/api/main.go