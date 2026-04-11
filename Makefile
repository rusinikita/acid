# Copy .env.example to .env
init:
	@if [ -f .env ]; then \
		echo "Error: .env already exists."; \
		exit 1; \
	else \
		cp .env.example .env; \
		echo ".env file created."; \
	fi

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# Run database in docker container
rundbs:
	docker-compose up -d

# Run postgres in docker container
runpg:
	docker-compose up -d postgres

# Run mysql in docker container
runmysql:
	docker-compose up -d mysql

# Stop containers and removing volumes
stopdb:
	docker-compose down -v

run:
	go run main.go

# Create and push a new tag: make tag VERSION=v0.1.0
tag:
	git tag $(VERSION)
	git push origin $(VERSION)

# Build and release with goreleaser (requires GITHUB_TOKEN)
release:
	goreleaser release --clean