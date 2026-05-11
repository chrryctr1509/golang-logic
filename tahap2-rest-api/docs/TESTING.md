# API Testing Documentation — tahap2-rest-api

## Setup

1. **Start server:**
   ```bash
   cd tahap2-rest-api
   go run main.go
   # Server running at http://localhost:8080
   ```

2. **Base URL:** `http://localhost:8080`

3. **Headers:** For authenticated endpoints, add:
   ```
   Authorization: Bearer {access_token}
   Content-Type: application/json
   ```

---

## 10 Test Methods

### TC-01: Register User (Success)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/register` |
| **Expected** | PASS — Status 201, response contains user_id |
| **Criteria** | Status 201, `"status": "SUCCESS"`, user_id returned |

**Request Body:**
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "0811999901",
    "address": "Jl. Sudirman No. 10",
    "pin": "123456"
}
```

**Response (201):**
```json
{
    "status": "SUCCESS",
    "result": {
        "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "first_name": "John",
        "last_name": "Doe",
        "phone_number": "0811999901",
        "address": "Jl. Sudirman No. 10",
        "balance": 0,
        "created_date": "2026-05-11T22:00:00+07:00"
    }
}
```

---

### TC-02: Register User (Duplicate Phone — FAIL Test)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/register` |
| **Expected** | FAIL Test — should return 400 error |
| **Criteria** | Status 400, `"Phone Number already registered"` |

**Request Body:**
```json
{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "0811999901",
    "address": "Jl. Sudirman No. 10",
    "pin": "123456"
}
```

**Response (400):**
```json
{
    "message": "Phone Number already registered"
}
```

---

### TC-03: Login (Success)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/login` |
| **Expected** | PASS — returns access_token + refresh_token |
| **Criteria** | Status 200, both tokens present, JWT format |

**Request Body:**
```json
{
    "phone_number": "0811999901",
    "pin": "123456"
}
```

**Response (200):**
```json
{
    "status": "SUCCESS",
    "result": {
        "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
        "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
    }
}
```

> Copy `access_token` for authenticated endpoints.

---

### TC-04: Login (Wrong PIN — FAIL Test)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/login` |
| **Expected** | FAIL Test — should return 400 error |
| **Criteria** | Status 400, `"Phone Number and PIN doesn't match."` |

**Request Body:**
```json
{
    "phone_number": "0811999901",
    "pin": "000000"
}
```

**Response (400):**
```json
{
    "message": "Phone Number and PIN doesn't match."
}
```

---

### TC-05: Top Up Balance (Success)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/api/v1/topup` |
| **Auth** | Bearer Token (required) |
| **Expected** | PASS — balance increases |
| **Criteria** | Status 200, `top_up_id` present, balance_after > balance_before |

**Request Body:**
```json
{
    "amount": 500000
}
```

**Response (200):**
```json
{
    "status": "SUCCESS",
    "result": {
        "top_up_id": "201ddde1-f797-484b-b1a0-07d1190e790a",
        "amount_top_up": 500000,
        "balance_before": 0,
        "balance_after": 500000,
        "created_date": "2026-05-11T22:10:00+07:00"
    }
}
```

---

### TC-06: Payment (Insufficient Balance — FAIL Test)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/api/v1/pay` |
| **Auth** | Bearer Token (required) |
| **Expected** | FAIL Test — should return 400 error |
| **Criteria** | Status 400, `"Balance is not enough"` |

**Request Body:**
```json
{
    "amount": 999999999,
    "remarks": "试图购买超出余额的商品"
}
```

**Response (400):**
```json
{
    "message": "Balance is not enough"
}
```

---

### TC-07: Payment (Success)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/api/v1/pay` |
| **Auth** | Bearer Token (required) |
| **Expected** | PASS — balance decreases |
| **Criteria** | Status 200, `payment_id` present, balance_after < balance_before |

**Request Body:**
```json
{
    "amount": 50000,
    "remarks": "购买电话卡 50 元"
}
```

**Response (200):**
```json
{
    "status": "SUCCESS",
    "result": {
        "payment_id": "13bcb11c-111e-4a65-9afd-90a86a01cd21",
        "amount": 50000,
        "remarks": "购买电话卡 50 元",
        "balance_before": 500000,
        "balance_after": 450000,
        "created_date": "2026-05-11T22:11:00+07:00"
    }
}
```

---

### TC-08: Transfer (Async — PASS Test)

| Field | Value |
|-------|-------|
| **Method** | POST |
| **URL** | `http://localhost:8080/api/v1/transfer` |
| **Auth** | Bearer Token (required) |
| **Expected** | PASS — immediate PENDING, worker changes to SUCCESS |
| **Criteria** | Immediate: `status: "PENDING"`. After 3s: DB shows `SUCCESS` |
| **Note** | Requires receiver `target_user` UUID from seed data |

**Request Body:**
```json
{
    "target_user": "b2c3d4e5-f6a7-8901-bcde-f12345678012",
    "amount": 30000,
    "remarks": "生日礼物"
}
```

**Immediate Response (200):**
```json
{
    "status": "SUCCESS",
    "result": {
        "transfer_id": "a7d39cf6-44b6-41fc-b3e9-7b16df5321c5",
        "amount": 30000,
        "remarks": "生日礼物",
        "balance_before": 450000,
        "balance_after": 450000,
        "created_date": "2026-05-11T22:12:00+07:00",
        "status": "PENDING"
    }
}
```

**DB Check (after 3 seconds):**
```sql
SELECT status, amount FROM transactions WHERE transaction_kind='TRANSFER' ORDER BY created_date DESC LIMIT 1;
-- Result: SUCCESS, 30000
```

---

### TC-09: Get Transactions (Dynamic Field Names)

| Field | Value |
|-------|-------|
| **Method** | GET |
| **URL** | `http://localhost:8080/api/v1/transactions` |
| **Auth** | Bearer Token (required) |
| **Expected** | PASS — field names vary by transaction kind |
| **Criteria** | TOPUP→`top_up_id`, PAYMENT→`payment_id`, TRANSFER→`transfer_id` |

**Response (200):**
```json
{
    "status": "SUCCESS",
    "result": [
        {
            "transfer_id": "a7d39cf6-44b6-41fc-b3e9-7b16df5321c5",
            "status": "SUCCESS",
            "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "transaction_type": "DEBIT",
            "amount": 30000,
            "remarks": "生日礼物",
            "balance_before": 450000,
            "balance_after": 420000,
            "created_date": "2026-05-11T22:12:00+07:00"
        },
        {
            "payment_id": "13bcb11c-111e-4a65-9afd-90a86a01cd21",
            "status": "SUCCESS",
            "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "transaction_type": "DEBIT",
            "amount": 50000,
            "remarks": "购买电话卡 50 元",
            "balance_before": 450000,
            "balance_after": 420000,
            "created_date": "2026-05-11T22:11:00+07:00"
        },
        {
            "top_up_id": "201ddde1-f797-484b-b1a0-07d1190e790a",
            "status": "SUCCESS",
            "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
            "transaction_type": "CREDIT",
            "amount": 500000,
            "remarks": "",
            "balance_before": 0,
            "balance_after": 500000,
            "created_date": "2026-05-11T22:10:00+07:00"
        }
    ]
}
```

---

### TC-10: Update Profile (Success)

| Field | Value |
|-------|-------|
| **Method** | PUT |
| **URL** | `http://localhost:8080/api/v1/profile` |
| **Auth** | Bearer Token (required) |
| **Expected** | PASS — name/address updated, phone unchanged |
| **Criteria** | Status 200, first_name/last_name/address changed, phone_number unchanged |

**Request Body:**
```json
{
    "first_name": "Tom",
    "last_name": "Araya",
    "address": "Jl. Diponegoro No. 215"
}
```

**Response (200):**
```json
{
    "status": "SUCCESS",
    "result": {
        "user_id": "a1b2c3d4-e5f6-7890-abcd-ef1234567890",
        "first_name": "Tom",
        "last_name": "Araya",
        "address": "Jl. Diponegoro No. 215",
        "phone_number": "0811999901",
        "balance": 420000,
        "updated_date": "2026-05-11T22:13:00+07:00"
    }
}
```

---

## Summary Table

| TC | Endpoint | Auth | Expected | PASS Criteria |
|----|----------|------|----------|----------------|
| TC-01 | POST /register | No | 201 SUCCESS | user_id returned |
| TC-02 | POST /register (dup) | No | 400 ERROR | "Phone Number already registered" |
| TC-03 | POST /login | No | 200 SUCCESS | access + refresh token |
| TC-04 | POST /login (wrong PIN) | No | 400 ERROR | "Phone Number and PIN doesn't match." |
| TC-05 | POST /topup | Bearer | 200 SUCCESS | top_up_id, balance_after > before |
| TC-06 | POST /pay (over balance) | Bearer | 400 ERROR | "Balance is not enough" |
| TC-07 | POST /pay | Bearer | 200 SUCCESS | payment_id, balance_after < before |
| TC-08 | POST /transfer | Bearer | 200 → DB SUCCESS | status PENDING immediately, SUCCESS in 3s |
| TC-09 | GET /transactions | Bearer | 200 SUCCESS | top_up_id/payment_id/transfer_id by kind |
| TC-10 | PUT /profile | Bearer | 200 SUCCESS | name changed, phone unchanged |