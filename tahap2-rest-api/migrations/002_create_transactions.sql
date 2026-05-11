-- Migration: Create transactions table
-- Run with: SET FOREIGN_KEY_CHECKS=0 at top, =1 at end

SET FOREIGN_KEY_CHECKS=0;

CREATE TABLE IF NOT EXISTS `transactions` (
    `transaction_id` CHAR(36) PRIMARY KEY,
    `user_id` CHAR(36) NOT NULL,
    `transaction_type` VARCHAR(10) NOT NULL,
    `transaction_kind` VARCHAR(20) NOT NULL,
    `amount` BIGINT NOT NULL,
    `remarks` TEXT,
    `balance_before` BIGINT NOT NULL,
    `balance_after` BIGINT NOT NULL,
    `status` VARCHAR(20) NOT NULL DEFAULT 'SUCCESS',
    `related_user_id` CHAR(36),
    `created_date` DATETIME NOT NULL,
    CONSTRAINT `fk_transactions_user` FOREIGN KEY (`user_id`) REFERENCES `users`(`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS=1;