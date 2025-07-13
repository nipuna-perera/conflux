# Conflux - AI Coding Agent Instructions

## Architecture Overview

**Conflux** is a configuration management platform with clean architecture separation:
- **Backend**: Go with repository pattern (`internal/repository/{mysql,postgres}`)
- **Frontend**: SvelteKit with TypeScript and Tailwind CSS
- **Database**: Dual support for MySQL and PostgreSQL via factory pattern
- **Auth**: JWT-based with middleware pipeline in `internal/api/middleware/`

## Development Workflow

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

## Integration Points

- **API Base URL**: Frontend uses `API_URL=http://backend:8080` in Docker
- **Authentication**: JWT tokens in Authorization header
- **CORS**: Configured in `main.go` with `handlers.CORS()`
- **Database migrations**: Auto-run on startup, manual via `make migrate`

When adding features, follow the layered architecture: models → repository interface → service → handler → routes.
