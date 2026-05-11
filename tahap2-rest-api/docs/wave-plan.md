# Wave Plan
generated_at: 2026-05-11T14:05:00Z
scope: NEW_FEATURE
effort: M

## Wave 1 — Foundation + Auth (Wave 1/2)

### Scope
Setup project structure, config, migrations, models, and auth endpoints.
Independent from Wave 2.

### Files
| Feature | Files | Agent |
|---------|-------|-------|
| go.mod + deps | go.mod, go.sum | be-developer |
| Config + DB | config/config.go, .env.example | be-developer |
| Migrations | migrations/001_create_users.sql, 002_create_transactions.sql, cmd/migrate/main.go | be-developer |
| Models | internal/model/user.go, transaction.go | be-developer |
| Repositories | internal/repository/user_repository.go, transaction_repository.go | be-developer |
| Auth Service | internal/service/auth_service.go, auth_service_test.go | be-developer |
| Auth Handler | internal/handler/auth_handler.go, internal/middleware/auth.go | be-developer |
| Main + Router | main.go | be-developer |

### Test Plan
- `go build ./...` passes
- Migration runs without error
- Auth endpoints respond (integration check via QA)

---

## Wave 2 — Wallet Operations + Worker (Wave 2/2)

### Scope
TopUp, Payment, Transfer (with background worker), Transactions Report, Profile Update.

### Files
| Feature | Files | Agent |
|---------|-------|-------|
| Transaction Service | internal/service/transaction_service.go, transaction_service_test.go | be-developer |
| Transaction Handler | internal/handler/transaction_handler.go, user_handler.go | be-developer |
| Transfer Worker | internal/worker/transfer_worker.go | be-developer |
| Seed Data | example_data/seed.sql | be-developer |
| README | README.md | be-developer |

### Test Plan
- `go test ./...` all PASS
- QA: test all 7 endpoints with curl/Postman

---

## Summary
| Wave | Features | Agent | Files |
|------|----------|-------|-------|
| 1 | Foundation + Auth | be-developer | ~12 |
| 2 | Wallet Ops + Worker | be-developer | ~8 |
