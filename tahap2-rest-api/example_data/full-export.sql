-- ============================================================
--  FULL EXPORT — tahap2-rest-api (wallet_db)
--  Generated: 2026-05-11
--  Contains: CREATE DATABASE + TABLES + SEED DATA
-- ============================================================

SET FOREIGN_KEY_CHECKS=0;

-- -----------------------------------------------------------
--  1. CREATE DATABASE
-- -----------------------------------------------------------
CREATE DATABASE IF NOT EXISTS wallet_db
    CHARACTER SET utf8mb4
    COLLATE utf8mb4_unicode_ci;

USE wallet_db;

-- -----------------------------------------------------------
--  2. CREATE TABLES
-- -----------------------------------------------------------

DROP TABLE IF EXISTS transactions;
DROP TABLE IF EXISTS users;

CREATE TABLE users (
    user_id       CHAR(36) NOT NULL,
    first_name    VARCHAR(100) NOT NULL,
    last_name     VARCHAR(100) NOT NULL,
    phone_number  VARCHAR(20) NOT NULL,
    address       TEXT,
    pin           VARCHAR(255) NOT NULL,
    balance       BIGINT NOT NULL DEFAULT 0,
    created_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_date  DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id),
    UNIQUE KEY phone_number (phone_number)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE transactions (
    transaction_id    CHAR(36) NOT NULL,
    user_id          CHAR(36) NOT NULL,
    transaction_type VARCHAR(10) NOT NULL,
    transaction_kind VARCHAR(20) NOT NULL,
    amount           BIGINT NOT NULL,
    remarks          TEXT,
    balance_before   BIGINT NOT NULL,
    balance_after    BIGINT NOT NULL,
    status           VARCHAR(20) NOT NULL DEFAULT 'SUCCESS',
    related_user_id  CHAR(36) DEFAULT NULL,
    created_date     DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (transaction_id),
    KEY fk_transactions_user (user_id),
    CONSTRAINT fk_transactions_user
        FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- -----------------------------------------------------------
--  3. SEED DATA — USERS
--  PIN for all seed users: "123456"
--  bcrypt hash: $2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe
-- -----------------------------------------------------------

INSERT INTO users (user_id, first_name, last_name, phone_number, address, pin, balance, created_date) VALUES
('a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'Alice',   'Wijaya',   '081234567891', 'Jl. Sudirman No. 10, Jakarta',   '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 500000, '2026-05-11 21:51:56'),
('b2c3d4e5-f6a7-8901-bcde-f12345678012', 'Bob',     'Santoso',  '081234567892', 'Jl. Gatot Subroto No. 25, Bandung', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 350000, '2026-05-11 21:51:56'),
('c3d4e5f6-a7b8-9012-cdef-012345678013', 'Charlie', 'Kusuma',   '081234567893', 'Jl. Ahmad Yani No. 88, Surabaya',  '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 100000, '2026-05-11 21:51:56');

-- -----------------------------------------------------------
--  4. SEED DATA — TRANSACTIONS
-- -----------------------------------------------------------

INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES

-- Alice transactions
('tx-alice-topup-001',    'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'CREDIT', 'TOPUP',     200000, 'Initial top-up',     0,      200000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-alice-topup-002',    'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'CREDIT', 'TOPUP',     500000, 'Salary deposit',     50000,  550000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-alice-pay-001',      'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'DEBIT',  'PAYMENT',   50000,  'Coffee shop',       200000, 150000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-alice-transfer-001', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'DEBIT',  'TRANSFER',  100000, 'Monthly allowance',  150000,  50000, 'PENDING', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', '2026-05-11 21:51:57'),

-- Bob transactions
('tx-bob-topup-001',      'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'CREDIT', 'TOPUP',     150000, 'Initial top-up',     0,      150000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-bob-pay-001',        'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'DEBIT',  'PAYMENT',   80000,  'Groceries',          250000, 170000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-bob-received-001',   'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'CREDIT', 'TRANSFER',  100000, 'From Alice',         150000, 250000, 'PENDING', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', '2026-05-11 21:51:57'),
('tx-bob-transfer-001',   'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'DEBIT',  'TRANSFER',  70000,  'Dinner split',        170000, 100000, 'PENDING', 'c3d4e5f6-a7b8-9012-cdef-012345678013', '2026-05-11 21:51:57'),

-- Charlie transactions
('tx-charlie-topup-001',  'c3d4e5f6-a7b8-9012-cdef-012345678013', 'CREDIT', 'TOPUP',     100000, 'Initial top-up',     0,      100000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-charlie-pay-001',    'c3d4e5f6-a7b8-9012-cdef-012345678013', 'DEBIT',  'PAYMENT',   20000,  'Online purchase',    170000, 150000, 'SUCCESS', NULL,                                  '2026-05-11 21:51:57'),
('tx-charlie-received-001','c3d4e5f6-a7b8-9012-cdef-012345678013', 'CREDIT', 'TRANSFER',  70000,  'From Bob',            100000, 170000, 'PENDING', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', '2026-05-11 21:51:57');

-- -----------------------------------------------------------
SET FOREIGN_KEY_CHECKS=1;

-- ============================================================
--  HOW TO USE
-- ============================================================
--
--  Option 1: Import entire file
--  $ mysql -u root -p123456 < full-export.sql
--
--  Option 2: Run line by line
--  $ mysql -u root -p123456
--  mysql> SOURCE full-export.sql;
--
--  Seed users login credentials:
--  +-----------------+------+----------+
--  | phone_number    | PIN  | balance  |
--  +-----------------+------+----------+
--  | 081234567891   | 123456 | 500000 |
--  | 081234567892   | 123456 | 350000 |
--  | 081234567893   | 123456 | 100000 |
--  +-----------------+------+----------+
--
-- ============================================================
