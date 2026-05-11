# Project Signal
generated_at: 2026-05-11T14:05:00Z

## Codebase Analysis
Repo: GREENFIELD — new project in tahap2-rest-api folder
No existing code. Pure new implementation.

Stack:
- Language: Go 1.21
- Framework: Gin (github.com/gin-gonic/gin)
- Database: MySQL 8.0+
- ORM: GORM v2
- Auth: JWT (github.com/golang-jwt/jwt/v5)
- Config: godotenv

## Scope
REST API Dompet Digital (e-wallet) dengan:
- User registration + PIN
- Login with JWT (access + refresh token)
- Top Up
- Payment
- Transfer (async via background worker + channel)
- Transaction history
- Profile update

## Files Expected
- go.mod, main.go, .env.example, README.md
- config/config.go
- cmd/migrate/main.go
- migrations/001_create_users.sql, 002_create_transactions.sql
- internal/middleware/auth.go
- internal/model/user.go, transaction.go
- internal/repository/user_repository.go, transaction_repository.go
- internal/service/auth_service.go, user_service.go, transaction_service.go
- internal/handler/auth_handler.go, user_handler.go, transaction_handler.go
- internal/worker/transfer_worker.go
- example_data/seed.sql
- Unit tests (bonus)

## Estimate
effort: M
risk_flags: background worker atomicity, JWT refresh token rotation
