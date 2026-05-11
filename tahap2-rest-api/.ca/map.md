# Project Map — tahap2-rest-api
generated_at: 2026-05-11T14:10:00Z

## Root Files
- go.mod / go.sum
- main.go                 → entry point, Gin router, inject deps, start worker
- .env.example
- README.md
- example_data/seed.sql

## cmd/
- migrate/main.go         → run SQL migrations

## config/
- config.go               → load .env, DB connection pool via GORM

## internal/middleware/
- auth.go                  → JWT extraction from Bearer token, inject user_id to gin.Context

## internal/model/
- user.go                  → User struct + GORM tags (CHAR(36) UUID, bcrypt PIN, balance BIGINT)
- transaction.go          → Transaction struct + enum constants (CREDIT/DEBIT, TOPUP/PAYMENT/TRANSFER, PENDING/SUCCESS/FAILED)

## internal/repository/
- user_repository.go      → CRUD user by phone, user by ID
- transaction_repository.go → Create tx, get user tx history, update tx status

## internal/service/
- auth_service.go          → Register (bcrypt PIN), Login (verify PIN, issue JWT), Refresh token
- user_service.go         → Update profile
- transaction_service.go   → TopUp, Payment (balance check + DB tx), Enqueue transfer job

## internal/handler/
- auth_handler.go          → POST /register, POST /login
- user_handler.go          → PUT /profile
- transaction_handler.go   → POST /topup, POST /pay, POST /transfer, GET /transactions

## internal/worker/
- transfer_worker.go       → Chan-based queue (TransferJob), Start(ctx), processTransfer (debit sender + credit receiver + DB tx)

## migrations/
- 001_create_users.sql
- 002_create_transactions.sql

## Tests
- internal/service/auth_service_test.go
- internal/service/transaction_service_test.go

## Stack
- Go 1.21 | Gin | GORM v2 | MySQL 8 | JWT | godotenv | bcrypt