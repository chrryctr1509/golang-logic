package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/user/tahap2-rest-api/internal/model"

	"gorm.io/gorm"
)

var ErrTransactionNotFound = errors.New("transaction not found")

type TransactionRepository interface {
	Create(ctx context.Context, tx *model.Transaction) error
	FindByUserID(ctx context.Context, userID string) ([]model.Transaction, error)
	FindByID(ctx context.Context, txID string) (*model.Transaction, error)
	UpdateStatus(ctx context.Context, txID string, status string) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	if err := r.db.WithContext(ctx).Create(tx).Error; err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}
	return nil
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID string) ([]model.Transaction, error) {
	var txns []model.Transaction
	if err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("created_date DESC").
		Find(&txns).Error; err != nil {
		return nil, fmt.Errorf("find transactions by user id: %w", err)
	}
	return txns, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, txID string) (*model.Transaction, error) {
	var tx model.Transaction
	if err := r.db.WithContext(ctx).Where("transaction_id = ?", txID).First(&tx).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTransactionNotFound
		}
		return nil, fmt.Errorf("find transaction by id: %w", err)
	}
	return &tx, nil
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, txID string, status string) error {
	if err := r.db.WithContext(ctx).
		Model(&model.Transaction{}).
		Where("transaction_id = ?", txID).
		Update("status", status).Error; err != nil {
		return fmt.Errorf("update transaction status: %w", err)
	}
	return nil
}
