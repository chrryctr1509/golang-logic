-- Migration: Create users table
-- Run with: SET FOREIGN_KEY_CHECKS=0 at top, =1 at end

SET FOREIGN_KEY_CHECKS=0;

CREATE TABLE IF NOT EXISTS `users` (
    `user_id` CHAR(36) PRIMARY KEY,
    `first_name` VARCHAR(100) NOT NULL,
    `last_name` VARCHAR(100) NOT NULL,
    `phone_number` VARCHAR(20) NOT NULL UNIQUE,
    `address` TEXT,
    `pin` VARCHAR(255) NOT NULL,
    `balance` BIGINT NOT NULL DEFAULT 0,
    `created_date` DATETIME NOT NULL,
    `updated_date` DATETIME NOT NULL ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

SET FOREIGN_KEY_CHECKS=1;