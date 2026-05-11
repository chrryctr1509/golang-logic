package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/user/tahap2-rest-api/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// Load config before running migrations
	_ = config.Load()

	if err := run(); err != nil {
		fmt.Println("Migration failed:", err)
		os.Exit(1)
	}
	fmt.Println("Migrations completed successfully.")
}

func run() error {
	cfg := config.GetConfig()

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.DBUser,
		cfg.DBPass,
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("get sql.DB: %w", err)
	}

	migrationFiles := []string{
		"migrations/001_create_users.sql",
		"migrations/002_create_transactions.sql",
	}

	for _, file := range migrationFiles {
		content, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read migration file %s: %w", file, err)
		}

		sanitized := removeFKCheck(content)

		fmt.Printf("Running migration: %s\n", file)
		if _, err := sqlDB.Exec(sanitized); err != nil {
			return fmt.Errorf("execute migration %s: %w", file, err)
		}
	}

	return nil
}

func removeFKCheck(content []byte) string {
	lines := strings.Split(string(content), "\n")
	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "SET FOREIGN_KEY_CHECKS") {
			continue
		}
		out = append(out, line)
	}
	return strings.Join(out, "\n")
}
