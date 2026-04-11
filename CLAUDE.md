# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run the app (standalone TUI)
make run
# or
go run main.go

# Start databases in Docker
make rundbs          # both PostgreSQL and MySQL
make runpg           # PostgreSQL only
make runmysql        # MySQL only
make stopdb          # stop and remove volumes

# Set up .env from example
make init

# Run all tests (uses testcontainers — Docker required)
go test ./...

# Run a single test file/package
go test ./tests/...
go test ./tests/ -run TestPostgresSuite/TestLostUpdate

# Run with a specific TOML sequence file
go run main.go -f path/to/sequence.toml
```

**Environment**: Copy `.env.example` to `.env`. Set `DB_DRIVER` (`postgres` or `mysql`) and `DB_CONNECT` (connection string).

## Architecture

ACID is a terminal-based tool for visualizing concurrent SQL transaction behavior. It has two operating modes:

### Standalone mode (`acid` / `acid -f file.toml`)
Single process: runs sequences against a local DB and shows results in a Bubbletea TUI.

### Client/server mode
- **`acid serve`** — starts the display TUI (HTTP server on `:7331`) with no DB connection; just renders events it receives
- **`acid run <name>` / `acid run -f file.toml`** — connects to a DB, executes the sequence, and streams events over HTTP to the server
- **`acid toggle`** — toggles result visibility on the running server (for quiz/prediction mode)
- **`acid init [dir]`** — scaffolds a learning environment (`.env`, `AGENTS.md`, `learning_plan.md`, example TOML sequences, Makefile)

### Package layout

| Package | Role |
|---|---|
| `call/` | DSL primitives: `Setup`, `Call`, `Begin`, `Commit`, `Rollback`. `Step` struct carries SQL code, transaction ID, and command type. |
| `sequence/` | `Sequence` struct (name, description, drop_tables, steps). `sequences.go` holds all built-in scenarios. |
| `runner/` | Executes a `Sequence` against a `*sql.DB`. Each step with a `TrxID` runs in its own goroutine with 200ms stagger. Emits `event.Event` values over a channel. `Iterator` wraps the channel for synchronous consumption. |
| `event/` | Immutable event values: `Call` (step started), `Result` (step finished), `Start`, `Done`. Carries waiting-transaction IDs for lock visualization. |
| `protocol/` | JSON wire format (`EventMessage`) for marshaling/unmarshaling events over HTTP between client and server. |
| `client/` | HTTP client that runs a sequence locally and POSTs each event to `acid serve`. |
| `server/` | HTTP server (`/event`, `/start`, `/done`, `/toggle-mode`, `/health`) that feeds events into a channel consumed by the TUI. |
| `config/` | Loads TOML sequence files into `sequence.Sequence`. |
| `db/` | Opens `*sql.DB` from env vars; loads `.env` files. |
| `ui/` | Bubbletea TUI: `router` (screen switching), `list` (sequence selector), `run` (events table with lock indicators and code highlighting), `theme`, `code`. |
| `initcmd/` | Scaffolds learning environments by copying embedded template files. |
| `terminal/` | Detects terminal type for split-pane instructions. |
| `tests/` | Integration tests using testcontainers (spins up `postgres:16-alpine`). |

### Data flow (standalone)

```
sequence.Sequence → runner.Run() → chan event.Event → runner.Iterator → ui/run (Bubbletea)
```

### Data flow (client/server)

```
sequence.Sequence → runner.Run() → chan event.Event → client.Run() → HTTP POST /event
                                                                            ↓
                                                               server.Server → chan event.Event → ui/run (Bubbletea)
```

## Adding a new built-in sequence

Add to `sequence/sequences.go` using the `call` DSL. Use `DropTables` to clean up before the sequence runs (preferred over `call.Setup("DROP TABLE IF EXISTS ...")` for tables that are reused across sequences).

## TOML sequence format

```toml
name = "Scenario Name"
description = "What this demonstrates"
drop_tables = ["table_name"]   # cleaned up before steps

[[steps]]
sql   = "CREATE TABLE ..."
setup = true   # hidden by default; press 's' to show

[[steps]]
cmd = "begin"   # begin | commit | rollback
trx = "alice"

[[steps]]
sql = "SELECT * FROM t"
trx = "alice"   # runs inside transaction "alice"

[[steps]]
sql = "SELECT * FROM t"
# no trx = auto-committed single query
```
