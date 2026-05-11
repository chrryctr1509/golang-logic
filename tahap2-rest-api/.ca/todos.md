# Todos — tahap2-rest-api
generated_at: 2026-05-11T14:10:00Z
status: active

## Wave 1 — Foundation + Auth
- [ ] go.mod + go.sum
- [ ] config/config.go + .env.example
- [ ] migrations/001_create_users.sql
- [ ] migrations/002_create_transactions.sql
- [ ] cmd/migrate/main.go
- [ ] internal/model/user.go
- [ ] internal/model/transaction.go
- [ ] internal/repository/user_repository.go
- [ ] internal/repository/transaction_repository.go
- [ ] internal/middleware/auth.go
- [ ] internal/service/auth_service.go
- [ ] internal/service/auth_service_test.go
- [ ] internal/handler/auth_handler.go
- [ ] main.go

## Wave 2 — Wallet Operations + Worker
- [ ] internal/service/transaction_service.go
- [ ] internal/service/transaction_service_test.go
- [ ] internal/worker/transfer_worker.go
- [ ] internal/handler/transaction_handler.go
- [ ] internal/handler/user_handler.go
- [ ] internal/service/user_service.go
- [ ] example_data/seed.sql
- [ ] README.md