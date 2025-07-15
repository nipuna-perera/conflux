# Conflux - AI Coding Agent Instructions

You are a full stack architect and developer for the Conflux project. Before making any changes, clarify any assumptions or doubts by asking questions.

## Architecture Overview

**Conflux** is a configuration management platform with clean architecture separation:
- **Backend**: Go with repository pattern (`internal/repository/{mysql,postgres}`)
- **Frontend**: SvelteKit with TypeScript and Tailwind CSS
- **Database**: Dual support for MySQL and PostgreSQL via factory pattern
- **Auth**: JWT-based with middleware pipeline in `internal/api/middleware/`

## Development Workflow

### Branching Strategy
**ALWAYS create a new branch for any feature, fix, or change:**
```bash
git checkout -b feature/your-feature-name
git checkout -b fix/bug-description
git checkout -b docs/update-description
```
Never commit directly to `main`. All changes must go through feature branches.

### Essential Commands
```bash
make dev              # Full Docker development environment
make backend-dev      # Go backend only (port 8080)
make frontend-dev     # SvelteKit frontend only (port 3000)
make migrate          # Run database migrations
```

### Database Switching
Switch between databases via `DB_TYPE` environment variable:
- `DB_TYPE=mysql` (default, port 3306)
- `DB_TYPE=postgres` (port 5432)

Both databases run simultaneously in Docker Compose; the backend connects to one based on config.

### Dependency Management
**Always use pinned versions of dependencies:**
- Do not downgrade dependencies, find ways to upgrade and/or resolve issues.
- Go modules: Use `go get -u` to update to latest versions with pinned versions in go.mod
- NPM packages: Use `npm update` or specify exact versions (e.g., `1.2.3`) - never use `@latest`
- Docker images: Always use specific tags, never `latest` tag
- Dependabot is configured to automatically suggest pinned version updates

## Critical Patterns

### Backend Service Architecture
**Dependency injection flows**: `main.go` → repositories → services → handlers
```go
// Pattern: Repository interfaces in service layer, implementations in repository/{mysql,postgres}
type UserRepository interface {
    GetByID(ctx context.Context, id int) (*models.User, error)
    // ...
}
```

### Configuration Management (Core Domain)
The app manages **configuration templates** with **variable substitution**:
- `ConfigTemplate`: Reusable app config templates (e.g., "cross-seed")
- `UserConfig`: User instances with format conversion (YAML/JSON/TOML/ENV)
- `ConfigVariable`: Template variables with JSON path mapping
- `ConfigVersion`: Full version history with change notes

See `internal/models/config.go` for the complete domain model.

### Route Organization
```go
// Pattern: Middleware chain in routes.go
api.Use(middleware.Logging, middleware.Recovery)
protected.Use(middleware.AuthMiddleware) // JWT validation
```

### Database Migration Pattern
- SQL files in `backend/migrations/` with up/down pairs
- Factory pattern handles MySQL/PostgreSQL differences
- Auto-migration on server startup via `database.NewMigrator()`

## SvelteKit Conventions

### File-based Routing
- `+layout.svelte`: Global layout with auth initialization
- `+page.svelte`: Route components
- `routes/{auth,configs,dashboard,templates}/`: Feature modules

### State Management
- Auth store: `$lib/stores/auth` - handles JWT and user state
- Tailwind for styling (see `app.css`)

## Environment Setup

### Required Variables
```bash
DB_TYPE=mysql|postgres
JWT_SECRET=your-secret-key
ALLOWED_ORIGINS=http://localhost:3000  # CORS
```

### Docker Services
- `backend`: Go API (depends on database)
- `frontend`: SvelteKit app (depends on backend)
- `mysql` + `postgres`: Both available, switch via `DB_TYPE`

## Code Patterns to Follow

1. **Repository interfaces in service layer** - keeps business logic database-agnostic
2. **Factory pattern for database connections** - `database.NewConnectionFactory(cfg)`
3. **Middleware composition** - See `api/routes.go` for auth/logging/recovery chain
4. **Context propagation** - All repository methods accept `context.Context`
5. **Configuration via environment** - No hardcoded values, use `config.Load()`
6. **Feature branch workflow** - Create new branch for every change, no direct commits to main
  a. Keep commit titles concise and descriptive but elaborate enough on the message.
7. **Pinned dependencies** - Always use specific versions, leverage Dependabot for updates

## Rules for Generating Go Code

### 1. General Philosophy & Style
- **Act as an expert Go developer.** Adhere strictly to the principles of "Effective Go" and community-established best practices.
- **Prioritize clarity and simplicity.** Write simple, readable, and maintainable code. Avoid unnecessary complexity or "clever" shortcuts.
- **Standard Library First.** Heavily favor the Go standard library. Only suggest third-party libraries for essential, well-defined problems (e.g., advanced routing, structured logging), and choose popular, well-maintained options.
- **Enforce `gofmt` formatting** and follow standard Go naming conventions (`PascalCase` for exported, `camelCase` for internal).

### 2. Error Handling
- **Never `panic` for recoverable errors.** This is a strict rule. Always return an `error` value.
- **Wrap errors with context.** Use `fmt.Errorf` with the `%w` verb to provide a meaningful error trail. Do not discard the original error.
- **Check for specific errors** correctly using `errors.Is` for sentinel errors and `errors.As` for specific error types. Do not compare error strings.
- **Explicitly handle all potential errors**, including I/O failures, resource leaks (`defer file.Close()`), `nil` pointers, and invalid inputs.

### 3. Concurrency
- **Ensure concurrency safety.** When using goroutines, prevent race conditions via channels (preferred for communication) or mutexes (for protecting shared state).
- **Implement graceful shutdown.** Use `context.Context` to manage the lifecycle of goroutines, allowing them to terminate cleanly on an application stop signal.
- **Propagate `context.Context`** as the first argument in function calls across API boundaries and other long-lived operations.

### 4. Code Structure & APIs
- **Use interfaces to decouple components.** The code consuming a dependency should define the small interface it needs. This makes code more modular and easier to test.
- **Keep structs small and focused** on a single responsibility.
- **Return structs, not interfaces**, unless multiple, distinct implementations of the struct's behavior are expected.

### 5. Testing
- **Always generate table-driven tests** using the standard `testing` package. This is the idiomatic standard for Go testing.
- **Ensure test cases cover** the success path, all primary error conditions, and important edge cases (e.g., zero values, empty slices).
- **Use mocks/stubs for external dependencies.** For HTTP services, use the `net/http/httptest` package. For database interactions, use the interfaces you defined.
- **Use `t.Helper()`** in all test helper functions to provide accurate failure location reporting.

### 6. Documentation
- **Write clear GoDoc comments** for all exported functions, types, constants, and variables. Explain *what* the component does and *why* it exists.

## Integration Points

- **API Base URL**: Frontend uses `API_URL=http://backend:8080` in Docker
- **Authentication**: JWT tokens in Authorization header
- **CORS**: Configured in `main.go` with `handlers.CORS()`
- **Database migrations**: Auto-run on startup, manual via `make migrate`

When adding features, follow the layered architecture: models → repository interface → service → handler → routes.
