-- Seed Data for tahap2-rest-api e-wallet
-- Run after migrations: mysql -u root -p wallet_db < example_data/seed.sql

SET FOREIGN_KEY_CHECKS=0;

-- Clear existing data
DELETE FROM transactions WHERE 1=1;
DELETE FROM users WHERE 1=1;

-- Insert 3 users with PIN "123456" (bcrypt hash)
-- PIN bcrypt hash for "123456": $2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe

INSERT INTO users (user_id, first_name, last_name, phone_number, address, pin, balance, created_date, updated_date) VALUES
('a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'Alice', 'Wijaya', '081234567891', 'Jl. Sudirman No. 10, Jakarta', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 500000, NOW(), NOW()),
('b2c3d4e5-f6a7-8901-bcde-f12345678012', 'Bob', 'Santoso', '081234567892', 'Jl. Gatot Subroto No. 25, Bandung', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 250000, NOW(), NOW()),
('c3d4e5f6-a7b8-9012-cdef-012345678013', 'Charlie', 'Kusuma', '081234567893', 'Jl. Ahmad Yani No. 88, Surabaya', '$2a$10$N9qo8uLOickgx2ZMRZoMy.MqrqJKyLpEqiFGJ9J0XwQxBJh1n.qe', 100000, NOW(), NOW());

-- Transactions for Alice (user_id: a1b2c3d4-e5f6-7890-abcd-ef1234567801)
-- TOPUP: 200000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-alice-topup-001', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'CREDIT', 'TOPUP', 200000, 'Initial top-up', 0, 200000, 'SUCCESS', NULL, NOW());

-- PAYMENT: 50000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-alice-pay-001', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'DEBIT', 'PAYMENT', 50000, 'Coffee shop', 200000, 150000, 'SUCCESS', NULL, NOW());

-- TRANSFER: 100000 to Bob
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-alice-transfer-001', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'DEBIT', 'TRANSFER', 100000, 'Monthly allowance', 150000, 50000, 'SUCCESS', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', NOW());

-- TOPUP: 500000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-alice-topup-002', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', 'CREDIT', 'TOPUP', 500000, 'Salary deposit', 50000, 550000, 'SUCCESS', NULL, NOW());

-- Transactions for Bob (user_id: b2c3d4e5-f6a7-8901-bcde-f12345678012)
-- TOPUP: 150000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-bob-topup-001', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'CREDIT', 'TOPUP', 150000, 'Initial top-up', 0, 150000, 'SUCCESS', NULL, NOW());

-- TRANSFER received from Alice: 100000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-bob-received-001', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'CREDIT', 'TRANSFER', 100000, 'From Alice', 150000, 250000, 'SUCCESS', 'a1b2c3d4-e5f6-7890-abcd-ef1234567801', NOW());

-- PAYMENT: 80000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-bob-pay-001', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'DEBIT', 'PAYMENT', 80000, 'Groceries', 250000, 170000, 'SUCCESS', NULL, NOW());

-- TRANSFER: 70000 to Charlie
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-bob-transfer-001', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', 'DEBIT', 'TRANSFER', 70000, 'Dinner split', 170000, 100000, 'SUCCESS', 'c3d4e5f6-a7b8-9012-cdef-012345678013', NOW());

-- Transactions for Charlie (user_id: c3d4e5f6-a7b8-9012-cdef-012345678013)
-- TOPUP: 100000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-charlie-topup-001', 'c3d4e5f6-a7b8-9012-cdef-012345678013', 'CREDIT', 'TOPUP', 100000, 'Initial top-up', 0, 100000, 'SUCCESS', NULL, NOW());

-- TRANSFER received from Bob: 70000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-charlie-received-001', 'c3d4e5f6-a7b8-9012-cdef-012345678013', 'CREDIT', 'TRANSFER', 70000, 'From Bob', 100000, 170000, 'SUCCESS', 'b2c3d4e5-f6a7-8901-bcde-f12345678012', NOW());

-- PAYMENT: 20000
INSERT INTO transactions (transaction_id, user_id, transaction_type, transaction_kind, amount, remarks, balance_before, balance_after, status, related_user_id, created_date) VALUES
('tx-charlie-pay-001', 'c3d4e5f6-a7b8-9012-cdef-012345678013', 'DEBIT', 'PAYMENT', 20000, 'Online purchase', 170000, 150000, 'SUCCESS', NULL, NOW());

SET FOREIGN_KEY_CHECKS=1;