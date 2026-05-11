# Dompet Digital — REST API E-Wallet

Go 1.21 | Gin | GORM v2 | MySQL 8 | JWT

## Setup

### 1. Clone & Install Dependencies

```bash
cd tahap2-rest-api
go mod tidy
```

### 2. Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

Required variables:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASS=yourpassword
DB_NAME=wallet_db
JWT_SECRET=your-secret-key-here
APP_PORT=8080
GIN_MODE=debug
```

### 3. Run Migrations

```bash
# Using the migration tool
go run cmd/migrate/main.go

# Or manually in MySQL
mysql -u root -p wallet_db < migrations/001_create_users.sql
mysql -u root -p wallet_db < migrations/002_create_transactions.sql
```

### 4. Seed Data (Optional)

```bash
mysql -u root -p wallet_db < example_data/seed.sql
```

Seed users (all PIN: `123456`):

| Name           | Phone        | Balance |
| -------------- | ------------ | ------- |
| Alice Wijaya   | 081234567891 | 500000  |
| Bob Santoso    | 081234567892 | 250000  |
| Charlie Kusuma | 081234567893 | 100000  |

### 5. Run Server

```bash
go run main.go
```

Server starts on `http://localhost:8080`.

---

## API Endpoints

### Public

#### Register

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "081299999999",
    "address": "Jl. Merdeka No. 123",
    "pin": "123456"
  }'
```

Response:

```json
{
  "status": "SUCCESS",
  "result": {
    "user_id": "uuid-here",
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "081299999999",
    "address": "Jl. Merdeka No. 123",
    "balance": 0,
    "created_date": "2026-05-11T00:00:00Z"
  }
}
```

#### Login

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "phone_number": "081299999999",
    "pin": "123456"
  }'
```

Response:

```json
{
  "status": "SUCCESS",
  "result": {
    "access_token": "eyJ...",
    "refresh_token": "eyJ..."
  }
}
```

---

### Protected (Bearer Token Required)

Get token from `/login`, then use in header:

```text
Authorization: Bearer <access_token>
```

#### Top Up

```bash
curl -X POST http://localhost:8080/api/v1/topup \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "amount": 100000
  }'
```

Response:

```json
{
  "status": "SUCCESS",
  "result": {
    "top_up_id": "uuid-here",
    "amount": 100000,
    "balance_before": 500000,
    "balance_after": 600000
  }
}
```

#### Payment

```bash
curl -X POST http://localhost:8080/api/v1/pay \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "amount": 50000,
    "remarks": "Coffee shop"
  }'
```

Response:

```json
{
  "status": "SUCCESS",
  "result": {
    "payment_id": "uuid-here",
    "amount": 50000,
    "remarks": "Coffee shop",
    "balance_before": 600000,
    "balance_after": 550000
  }
}
```

#### Transfer

```bash
curl -X POST http://localhost:8080/api/v1/transfer \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "target_user": "b2c3d4e5-f6a7-8901-bcde-f12345678012",
    "amount": 50000,
    "remarks": "Monthly allowance"
  }'
```

Response (async — status is PENDING, then SUCCESS/FAILED processed by worker):

```json
{
  "status": "SUCCESS",
  "result": {
    "transfer_id": "uuid-here",
    "amount": 50000,
    "remarks": "Monthly allowance",
    "balance_before": 550000,
    "balance_after": 550000,
    "status": "PENDING"
  }
}
```

#### Get Transactions

```bash
curl -X GET http://localhost:8080/api/v1/transactions \
  -H "Authorization: Bearer <access_token>"
```

Response (dynamic field name by transaction kind):

```json
{
  "status": "SUCCESS",
  "result": [
    {
      "top_up_id": "uuid-here",
      "amount": 100000,
      "balance_before": 0,
      "balance_after": 100000,
      "status": "SUCCESS",
      "created_date": "2026-05-11T00:00:00Z"
    },
    {
      "payment_id": "uuid-here",
      "amount": 50000,
      "remarks": "Coffee shop",
      "balance_before": 100000,
      "balance_after": 50000,
      "status": "SUCCESS",
      "created_date": "2026-05-11T00:00:01Z"
    },
    {
      "transfer_id": "uuid-here",
      "amount": 50000,
      "remarks": "Monthly allowance",
      "balance_before": 50000,
      "balance_after": 50000,
      "status": "SUCCESS",
      "related_user_id": "b2c3d4e5-f6a7-8901-bcde-f12345678012",
      "created_date": "2026-05-11T00:00:02Z"
    }
  ]
}
```

#### Update Profile

```bash
curl -X PUT http://localhost:8080/api/v1/profile \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <access_token>" \
  -d '{
    "first_name": "John",
    "last_name": "Smith",
    "address": "Jl. Baru No. 99, Jakarta"
  }'
```

Response:

```json
{
  "status": "SUCCESS",
  "result": {
    "user_id": "uuid-here",
    "first_name": "John",
    "last_name": "Smith",
    "phone_number": "081299999999",
    "address": "Jl. Baru No. 99, Jakarta",
    "balance": 550000,
    "updated_date": "2026-05-11T00:00:00Z"
  }
}
```

---

## Error Responses

| HTTP | Condition                | Response                                             |
| ---- | ------------------------ | ---------------------------------------------------- |
| 400  | Amount <= 0              | `{"message": "Amount must be greater than 0"}`       |
| 400  | Balance insufficient     | `{"message": "Balance is not enough"}`               |
| 400  | Phone already registered | `{"message": "Phone Number already registered"}`     |
| 400  | Wrong PIN                | `{"message": "Phone Number and PIN doesn't match."}` |
| 401  | Invalid/missing token    | `{"message": "Unauthenticated"}`                     |
| 404  | User not found           | `{"message": "User not found"}`                      |

---

## Running Tests

```bash
go test ./internal/service/...
```

---

## Architecture

```text
main.go
  └── Gin Router
        ├── Public routes: /register, /login
        └── Protected routes (/api/v1):
              ├── POST /topup       → TransactionHandler.TopUp
              ├── POST /pay         → TransactionHandler.Payment
              ├── POST /transfer    → TransactionHandler.Transfer
              ├── GET /transactions → TransactionHandler.GetTransactions
              └── PUT /profile      → UserHandler.UpdateProfile

Services:
  ├── AuthService        → Register, Login, JWT generation
  ├── TransactionService → TopUp, Payment, Transfer, GetTransactions
  └── UserService        → UpdateProfile

Worker:
  └── TransferWorker     → Channel-based async transfer processor
                          (row lock + DB transaction)

Repositories:
  ├── UserRepository        → CRUD users
  └── TransactionRepository → CRUD transactions
```
