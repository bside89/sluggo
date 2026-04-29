.DEFAULT_GOAL := help

.PHONY: help local prod docs build

help:
	@echo "Available commands:"
	@echo "  make local  	Start PostgreSQL and run the API locally"
	@echo "  make prod  	Build and start all services with Docker Compose"
	@echo "  make docs    	Regenerate Swagger documentation"
	@echo "  make build   	Build the API binary"

local:
	@bash ./scripts/run-local.sh

prod:
	@bash ./scripts/run-prod.sh

docs:
	@bash ./scripts/update-swagger.sh

build:
	@go build -o sluggo cmd/api/main.go