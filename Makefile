# Copy .env.example to .env
init:
	@if [ -f .env ]; then \
		echo "Error: .env already exists."; \
		exit 1; \
	else \
		cp .env.example .env; \
		echo ".env file created."; \
	fi

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