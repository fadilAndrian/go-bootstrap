# Go Bootstrap

A modular, extensible Go project generator CLI designed for production-grade applications.

Select the components you need---database, HTTP framework, background worker---and generate ordinary, idiomatic, editable Go source code with **no runtime lock-in**.

---

## Features

- **Modular Components:** Choose your HTTP router, database adapter, and background worker independently.
- **No Runtime Lock-In:** Generated projects are standard Go code with explicit lifecycle and configuration. No runtime CLI dependency.
- **Robust Lifecycle Management:** Coordinated startup, automatic rollback on initialization failure, and graceful shutdown (HTTP drain $\rightarrow$ worker drain $\rightarrow$ database close).
- **Compatibility Validation:** Automatically prevents incompatible component combinations (e.g., rejecting River background worker without PostgreSQL).
- **Project Manifest:** Tracks generator-managed components and versions in `.go-bootstrap.yaml`.

---

## Installation

```bash
go install github.com/konnco/go-bootstrap/cmd/go-bootstrap@latest
```

Or build from source:

```bash
git clone https://github.com/konnco/go-bootstrap.git
cd go-bootstrap
go build -o bin/go-bootstrap ./cmd/go-bootstrap
```

---

## Quick Start

### 1. Create a New Project

```bash
go-bootstrap new my-app \
  --db=postgres \
  --http=gin \
  --worker=river \
  --module=github.com/myuser/my-app
```

### 2. Run the Generated Application

```bash
cd my-app
cp .env.example .env
go-bootstrap serve
# or directly: go run ./cmd/serve
```

---

## Available Components

Run `go-bootstrap list` to inspect all supported plugins:

| Category | Plugin | Description |
| :--- | :--- | :--- |
| **database** | `postgres` | PostgreSQL connection pool with `pgx/v5` |
| | `mysql` | MySQL adapter using `go-sql-driver/mysql` |
| | `sqlite` | Pure Go SQLite adapter using `modernc.org/sqlite` |
| | `none` | No database |
| **http** | `gin` | Gin high-performance HTTP web framework |
| | `chi` | Chi lightweight and composable HTTP router |
| | `nethttp` | Standard library `net/http` ServeMux |
| **worker** | `river` | Fast, robust PostgreSQL-backed background job queue |
| | `none` | No background worker |

---

## CLI Commands

| Command | Usage | Description |
| :--- | :--- | :--- |
| `new` | `go-bootstrap new <project-name> [flags]` | Generates a new project in a new folder |
| `init` | `go-bootstrap init [flags]` | Initializes the current folder as a project |
| `list` | `go-bootstrap list` | Lists all available component plugins |
| `doctor` | `go-bootstrap doctor` | Validates project structure, manifest, and compatibility |
| `serve` | `go-bootstrap serve` | Convenience runner for `go run ./cmd/serve` |

---

## Generated Project Architecture

```text
my-project/
├── cmd/
│   └── serve/
│       └── main.go              # Composition root and OS signal handling
├── internal/
│   ├── bootstrap/
│   │   ├── app.go               # Startup, component wiring, graceful shutdown
│   │   └── lifecycle.go         # LIFO resource cleanup coordinator
│   ├── config/
│   │   └── config.go            # Structured config loader with env defaults
│   ├── database/
│   │   └── postgres.go          # Database connection pool and ping checks
│   ├── http/
│   │   ├── server.go            # HTTP server wrapper with bounded shutdown
│   │   └── gin.go               # Gin engine setup and middleware
│   ├── worker/
│   │   ├── river.go             # River client lifecycle management
│   │   └── job.go               # Sample background job and worker
│   └── routes/
│       └── routes.go            # API route registration
├── modules/
│   └── example/
│       └── handler.go           # Example domain handlers (/health, /api/v1/...)
├── configs/
│   └── config.yaml              # Default configuration file
├── .env.example                 # Example environment variables
├── .go-bootstrap.yaml           # Generator manifest
├── go.mod
└── go.sum
```

---

## License

MIT
