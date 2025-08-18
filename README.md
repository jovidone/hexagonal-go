# Hexagonal Go

This project showcases a simple banking API built in Go using the **Hexagonal Architecture** pattern. It demonstrates how to structure code around domain rules and ports/adapters while using [Gin](https://github.com/gin-gonic/gin) for HTTP and [GORM](https://gorm.io/) for database access.

## Project Structure
```
cmd/                 Application entry point
internal/
  adapters/          HTTP handlers and database adapters
  config/            Database configuration and migration
  core/              Domain, ports, and services
migrations/          Docker compose for local PostgreSQL
```

## Requirements
- Go 1.20 or later
- Docker (for running PostgreSQL locally)

## Setup
### 1. Configure Environment Variables
Copy `.env.example` to `.env` and adjust values for your environment:
```bash
cp .env.example .env
```
Required variables:
- `DB_HOST`
- `DB_USER`
- `DB_PASSWORD`
- `DB_NAME`
- `DB_PORT`
- `DB_SSLMODE` (defaults to `disable` if unset)
- `JWT_SECRET`

### 2. Start the Database (optional)
A docker-compose file is provided for local development:
```bash
cd migrations
docker compose up -d
cd ..
```
This starts PostgreSQL on port **5432** with default credentials defined in the compose file.

### 3. Install Go Dependencies
```bash
go mod tidy
```

## Running the Application
```bash
go run cmd/main.go
```
The server starts on port `8080` and automatically runs database migrations.

## API Endpoints
| Method | Path                         | Description                |
|--------|------------------------------|----------------------------|
| POST   | `/register`                  | Register a new user        |
| POST   | `/login`                     | Authenticate and receive tokens |
| POST   | `/refresh`                   | Refresh JWT token          |
| POST   | `/deposit`                   | Deposit funds *(auth required)* |
| POST   | `/withdraw`                  | Withdraw funds *(auth required)* |
| POST   | `/transfer`                  | Transfer funds *(auth required)* |
| GET    | `/transactions/:user_id`     | List user transactions *(auth required)* |
| GET    | `/profile`                   | Retrieve user profile *(auth required)* |

## Running Tests
Unit tests are provided for core services:
```bash
go test ./...
```

## Additional Notes
- The database connection enables the `uuid-ossp` extension and runs automatic migrations for the `User` and `Transaction` models.
- This repository is intended for learning and experimentation with the hexagonal architecture approach in Go.

