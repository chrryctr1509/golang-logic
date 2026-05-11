# Technical Spec — tahap2-rest-api
generated_at: 2026-05-11T14:10:00Z

## Overview
REST API Dompet Digital (e-wallet) dengan Clean Architecture.
Go 1.21 | Gin | GORM v2 | MySQL 8 | JWT | bcrypt

---

## 1. Project Structure

```
tahap2-rest-api/
├── go.mod
├── go.sum
├── .env.example
├── README.md
├── main.go                           # entry point
├── example_data/seed.sql
├── cmd/migrate/main.go               # migration runner
├── config/config.go                  # load env + DB
├── internal/
│   ├── middleware/auth.go            # JWT middleware
│   ├── model/
│   │   ├── user.go                   # User model
│   │   └── transaction.go            # Transaction model + enums
│   ├── repository/
│   │   ├── user_repository.go
│   │   └── transaction_repository.go
│   ├── service/
│   │   ├── auth_service.go
│   │   ├── user_service.go
│   │   ├── transaction_service.go
│   │   ├── auth_service_test.go
│   │   └── transaction_service_test.go
│   ├── handler/
│   │   ├── auth_handler.go
│   │   ├── user_handler.go
│   │   └── transaction_handler.go
│   └── worker/transfer_worker.go
└── migrations/
    ├── 001_create_users.sql
    └── 002_create_transactions.sql
```

---

## 2. Config (config/config.go)

Load dari `.env` via godotenv:
- DB_HOST, DB_PORT, DB_USER, DB_PASS, DB_NAME
- JWT_SECRET
- APP_PORT (default 8080)

DSN: `fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", ...)`

Expose: `GetDB() *gorm.DB`, `GetConfig() Config`

---

## 3. Models

### User
```go
type User struct {
    UserID      string    `gorm:"column:user_id;type:char(36);primaryKey"`
    FirstName   string    `gorm:"column:first_name;type:varchar(100)"`
    LastName    string    `gorm:"column:last_name;type:varchar(100)"`
    PhoneNumber string    `gorm:"column:phone_number;type:varchar(20);uniqueIndex"`
    Address     string    `gorm:"column:address;type:text"`
    PIN         string    `gorm:"column:pin;type:varchar(255)"`
    Balance     int64     `gorm:"column:balance;type:bigint;default:0"`
    CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
    UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}
```

### Transaction
```go
const (
    TypeCredit  = "CREDIT"
    TypeDebit   = "DEBIT"
    KindTopup   = "TOPUP"
    KindPayment = "PAYMENT"
    KindTransfer= "TRANSFER"
    StatusPending  = "PENDING"
    StatusSuccess  = "SUCCESS"
    StatusFailed   = "FAILED"
)

type Transaction struct {
    TransactionID   string    `gorm:"column:transaction_id;type:char(36);primaryKey"`
    UserID          string    `gorm:"column:user_id;type:char(36)"`
    TransactionType string    `gorm:"column:transaction_type;type:varchar(10)"`
    TransactionKind string    `gorm:"column:transaction_kind;type:varchar(20)"`
    Amount          int64     `gorm:"column:amount;type:bigint"`
    Remarks         string    `gorm:"column:remarks;type:text"`
    BalanceBefore   int64     `gorm:"column:balance_before;type:bigint"`
    BalanceAfter    int64     `gorm:"column:balance_after;type:bigint"`
    Status          string    `gorm:"column:status;type:varchar(20);default:SUCCESS"`
    RelatedUserID   string    `gorm:"column:related_user_id;type:char(36)"`
    CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime"`
}
```

---

## 4. Repository Interfaces

### UserRepository
- `Create(ctx, user *User) error`
- `FindByPhone(ctx, phone) (*User, error)`
- `FindByID(ctx, userID) (*User, error)`
- `Update(ctx, user *User) error`

### TransactionRepository
- `Create(ctx, tx *Transaction) error`
- `FindByUserID(ctx, userID) ([]Transaction, error)`
- `FindByID(ctx, txID) (*Transaction, error)`
- `UpdateStatus(ctx, txID, status) error`

---

## 5. JWT Spec

Token claims:
```go
type Claims struct {
    UserID      string `json:"user_id"`
    PhoneNumber string `json:"phone_number"`
    jwt.RegisteredClaims
}
```

- access_token: 15 min expiry (time.Now().Add(15 * time.Minute))
- refresh_token: 7 days expiry (time.Now().Add(7 * 24 * time.Hour))
- Algorithm: HS256
- Signing key: from JWT_SECRET env

---

## 6. Transfer Worker

```go
type TransferJob struct {
    TransactionID string
    SenderID     string
    ReceiverID   string
    Amount       int64
}

type TransferWorker struct {
    Queue chan TransferJob
    db    *gorm.DB
    // repositories injected
}

func (w *TransferWorker) Start(ctx context.Context) {
    for job := range w.Queue {
        w.processTransfer(job)
    }
}

func (w *TransferWorker) processTransfer(job TransferJob) {
    // BEGIN
    // SELECT FOR UPDATE sender balance
    // Check balance >= amount
    // Debit sender (balance -= amount)
    // Credit receiver (balance += amount)
    // Update tx status = SUCCESS
    // COMMIT
    // On error: ROLLBACK, update tx status = FAILED
}
```

Key: use `db.Transaction()` + `db.Clauses(clause.Locking{Strength: "UPDATE"})` for row-level lock.

---

## 7. API Response Format

Success:
```json
{
    "status": "SUCCESS",
    "result": { ... }
}
```

Error:
```json
{
    "message": "error text"
}
```

Dynamic field mapping in GET /transactions:
- kind==TOPUP → "top_up_id"
- kind==PAYMENT → "payment_id"
- kind==TRANSFER → "transfer_id"

---

## 8. Migration

- SQL files in `migrations/`
- cmd/migrate/main.go: read + execute all .sql files in order
- SET FOREIGN_KEY_CHECKS=0 at start, =1 at end
- Or use GORM AutoMigrate as fallback

---

## 9. Key Implementation Details

1. **UUID generation**: `uuid.New().String()` in Go, NOT DB-generated
2. **PIN hashing**: `bcrypt.GenerateFromPassword([]byte(pin), bcrypt.DefaultCost)`
3. **PIN verify**: `bcrypt.CompareHashAndPassword([]byte(hash), []byte(pin))`
4. **Balance update**: GORM transaction with row lock (`Clauses(Locking{Strength: "UPDATE"})`)
5. **Transfer**: immediate PENDING response, worker async, update to SUCCESS/FAILED
6. **Refresh token**: NOT implemented (out of scope per brief), just access_token

---

## 10. Dependencies (go.mod)

```
module github.com/user/tahap2-rest-api
go 1.21

require (
    github.com/gin-gonic/gin v1.9.1
    github.com/golang-jwt/jwt/v5 v5.2.0
    github.com/google/uuid v1.5.0
    gorm.io/gorm v1.25.5
    gorm.io/driver/mysql v1.5.2
    github.com/joho/godotenv v1.5.1
    golang.org/x/crypto v0.18.0
    github.com/stretchr/testify v1.8.4
)
```

---

## 11. File → Layer Mapping

| File | Layer | Responsibility |
|------|-------|--------------|
| config/config.go | Config | env loading, DB pool |
| model/*.go | Model | struct, constants, GORM tags |
| repository/*.go | Repository | DB queries, transactions |
| service/*.go | Service | business logic, validation |
| handler/*.go | Handler | HTTP request/response |
| middleware/auth.go | Middleware | JWT verify, user inject |
| worker/*.go | Worker | async transfer processing |
| main.go | Entry | router setup, DI, start worker |