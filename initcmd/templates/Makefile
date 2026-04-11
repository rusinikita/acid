# Start PostgreSQL (matches .env default)
pg:
	docker-compose up -d postgres

# Start MySQL (switch DB_DRIVER in .env)
mysql:
	docker-compose up -d mysql

# Stop all database containers
stop:
	docker-compose down

# Start acid server (default port 7331)
serve:
	acid serve

# Update acid — detects Homebrew or falls back to install script
update:
	@if command -v brew >/dev/null 2>&1 && brew list --cask acid >/dev/null 2>&1; then \
		brew upgrade --cask acid; \
	else \
		curl -fsSL https://raw.githubusercontent.com/rusinikita/acid/main/install.sh | sh; \
	fi

.PHONY: pg mysql stop serve update
