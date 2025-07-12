# Full-Stack Application

A modern full-stack web application built with Go backend and SvelteKit frontend.

## Architecture

- **Backend**: Go with clean architecture (cmd/internal structure)
- **Frontend**: SvelteKit with TypeScript
- **Database**: Support for both MySQL and PostgreSQL
- **Authentication**: JWT-based authentication
- **Containerization**: Docker and Docker Compose ready

## Project Structure

```
├── backend/              # Go backend application
│   ├── cmd/             # Application entry points
│   ├── internal/        # Private application code
│   └── pkg/             # Public packages
├── frontend/            # SvelteKit frontend application
│   └── src/             # Source code
├── scripts/             # Development scripts
└── docker-compose.yml   # Container orchestration
```

## Quick Start

1. **Setup environment**:
   ```bash
   chmod +x scripts/setup.sh
   ./scripts/setup.sh
   ```

2. **Start development environment**:
   ```bash
   make dev
   ```

3. **Access the application**:
   - Frontend: http://localhost:3000
   - Backend API: http://localhost:8080

## Development Commands

- `make dev` - Start full development environment
- `make backend-dev` - Run backend only
- `make frontend-dev` - Run frontend only
- `make test` - Run all tests
- `make migrate` - Run database migrations
- `make clean` - Clean up containers and volumes

## Environment Configuration

Copy `.env.example` to `.env` and configure your environment variables:

```bash
cp .env.example .env
```

## Database Support

The application supports both MySQL and PostgreSQL. Configure the database type in your environment variables:

- Set `DB_TYPE=mysql` for MySQL
- Set `DB_TYPE=postgres` for PostgreSQL

## API Documentation

The REST API provides the following endpoints:

- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/profile` - Update user profile

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

This project is licensed under the MIT License.
