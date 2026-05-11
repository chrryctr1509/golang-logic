# Brief Extract — tahap2-rest-api
generated_at: 2026-05-11T14:10:00Z
from: briefs/brief2.md

## Features
1. **Register** — POST /register (phone_number unique, PIN bcrypt, return user data)
2. **Login** — POST /login (verify PIN, return JWT access + refresh token)
3. **Top Up** — POST /topup (auth, add balance, record transaction)
4. **Payment** — POST /pay (auth, deduct balance, record transaction)
5. **Transfer** — POST /transfer (auth, PENDING → worker async, debit sender + credit receiver)
6. **Transactions Report** — GET /transactions (auth, all user tx, dynamic field names by kind)
7. **Update Profile** — PUT /profile (auth, update first_name/last_name/address, NOT phone/pin)

---

## API Endpoints

| Method | Path | Auth | Request | Response |
|--------|------|------|---------|---------|
| POST | /register | No | first_name, last_name, phone_number, address, pin | user data |
| POST | /login | No | phone_number, pin | access_token + refresh_token |
| POST | /topup | Bearer | amount | top_up_id, amount, balance_before, balance_after |
| POST | /pay | Bearer | amount, remarks | payment_id, amount, remarks, balance_before, balance_after |
| POST | /transfer | Bearer | target_user (UUID), amount, remarks | transfer_id, amount, remarks, balance_before, balance_after |
| GET | /transactions | Bearer | — | array of transactions |
| PUT | /profile | Bearer | first_name, last_name, address | updated user data |

---

## Database Schema

### users
```sql
user_id       CHAR(36) PRIMARY KEY  -- UUID from Go
first_name    VARCHAR(100)
last_name     VARCHAR(100)
phone_number  VARCHAR(20) UNIQUE
address       TEXT
pin           VARCHAR(255)          -- bcrypt hash
balance       BIGINT DEFAULT 0
created_date  DATETIME
updated_date  DATETIME ON UPDATE
```

### transactions
```sql
transaction_id    CHAR(36) PRIMARY KEY
user_id           CHAR(36) FK → users.user_id
transaction_type  VARCHAR(10)  -- CREDIT | DEBIT
transaction_kind  VARCHAR(20)  -- TOPUP | PAYMENT | TRANSFER
amount            BIGINT
remarks           TEXT
balance_before    BIGINT
balance_after     BIGINT
status            VARCHAR(20) DEFAULT SUCCESS  -- PENDING | SUCCESS | FAILED
related_user_id   CHAR(36)     -- receiver for TRANSFER
created_date      DATETIME
```

---

## Auth Spec

- **access_token**: expire 15 min, payload: user_id, phone_number
- **refresh_token**: expire 7 days
- Library: github.com/golang-jwt/jwt/v5
- Middleware: extract Bearer token, set user_id in gin.Context
- Error: 401 Unauthenticated

---

## Background Worker

- `/transfer` immediately returns PENDING response
- Goroutine worker runs via `go worker.Start(ctx)`
- Queue: Go channel `chan TransferJob`
- Process: BEGIN → debit sender → credit receiver → COMMIT → update tx status SUCCESS
- On error: rollback, update tx status FAILED

```go
type TransferJob struct {
    TransactionID string
    SenderID     string
    ReceiverID   string
    Amount       int64
}
```

---

## Error Handling

| Condition | HTTP | Response |
|-----------|------|---------|
| Token invalid/expired | 401 | {"message": "Unauthenticated"} |
| Balance insufficient | 400 | {"message": "Balance is not enough"} |
| Phone registered | 400 | {"message": "Phone Number already registered"} |
| Phone/PIN mismatch | 400 | {"message": "Phone Number and PIN doesn't match."} |
| Target user not found | 404 | {"message": "User not found"} |
| Amount <= 0 | 400 | {"message": "Amount must be greater than 0"} |

---

## Non-Functional
- MySQL: CHAR(36) for UUID, ENGINE=InnoDB, parseTime=True
- Balance: BIGINT (rupiah integer, no float)
- PIN: bcrypt hash only, never plain text
- DB transactions for balance updates (atomicity)
- Dynamic response field names by transaction_kind

---

## Unit Test (Bonus)
- auth_service_test.go: login success, login fail (wrong PIN), duplicate register
- transaction_service_test.go: topup, payment (sufficient/insufficient balance), transfer
- Library: github.com/stretchr/testify