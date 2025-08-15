# Hexagonal Go

This application demonstrates a basic Hexagonal architecture using Gin and Gorm.

## Setup

1. Copy `.env.example` to `.env` and adjust the values:
   ```bash
   cp .env.example .env
   ```
   Required variables:
   - `DB_HOST`
   - `DB_USER`
   - `DB_PASSWORD`
   - `DB_NAME`
   - `DB_PORT`
   - `DB_SSLMODE`
   - `JWT_SECRET`

2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/main.go
   ```

The server starts on port 8080.
