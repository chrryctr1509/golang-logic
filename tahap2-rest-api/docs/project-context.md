# Project Context — tahap2-rest-api
generated_at: 2026-05-11T14:10:00Z

## Overview
REST API Dompet Digital (e-wallet). Go 1.21 + Gin + MySQL + JWT. Clean Architecture.

## Stack
- Go 1.21, Gin framework, GORM v2, MySQL 8, JWT (golang-jwt/jwt/v5), godotenv, bcrypt

## Key Design Decisions
1. **UUID di aplikasi** — MySQL CHAR(36), `uuid.New().String()` di Go, bukan DB-generated
2. **Balance BIGINT** — simpan sebagai sen/rupiah integer, bukan float
3. **bcrypt PIN** — hash sebelum simpan, verify saat login
4. **JWT 2 token** — access (15 min) + refresh (7 days), payload: user_id, phone_number
5. **Transfer async** — endpoint return PENDING immediately, goroutine worker proses via channel
6. **DB Transaction** — balance update atomic via GORM Transaction()
7. **MySQL only** — ENGINE=InnoDB, FOREIGN KEY, parseTime=True

## Layer Pattern
- Handler: parse request, call service, format response
- Service: business logic, validation, enqueue worker jobs
- Repository: raw DB queries via GORM

## Project Structure
```
tahap2-rest-api/
├── main.go
├── config/config.go
├── cmd/migrate/main.go
├── internal/{middleware,model,repository,service,handler,worker}/
├── migrations/
├── example_data/seed.sql
└── .env.example
```

## From Brief
- 7 endpoints: /register, /login, /topup, /pay, /transfer, /transactions, /profile
- 2 JWT tokens, PIN bcrypt, background transfer worker
- Dynamic response field names (top_up_id / payment_id / transfer_id)
- seed.sql dengan 3 user + sample transactions

## Conventions Applied
- Go file: snake_case (user_repository.go)
- Struct: PascalCase (UserRepository)
- Interface: PascalCase, suffix -er (UserRepository interface)
- Error: wrap with fmt.Errorf("context: %w", err)
- Test: *_test.go, testify assertions