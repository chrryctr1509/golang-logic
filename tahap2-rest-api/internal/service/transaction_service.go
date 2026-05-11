package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/user/tahap2-rest-api/internal/model"
	"github.com/user/tahap2-rest-api/internal/repository"
	"github.com/user/tahap2-rest-api/internal/worker"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var (
	ErrInsufficientBalance = errors.New("balance is not enough")
	ErrInvalidAmount       = errors.New("amount must be greater than 0")
)

type TransactionService struct {
	db          *gorm.DB
	txRepo      repository.TransactionRepository
	userRepo    repository.UserRepository
	transferQ   chan worker.TransferJob
}

func NewTransactionService(db *gorm.DB, txRepo repository.TransactionRepository, userRepo repository.UserRepository, transferQ chan worker.TransferJob) *TransactionService {
	return &TransactionService{
		db:        db,
		txRepo:    txRepo,
		userRepo:  userRepo,
		transferQ: transferQ,
	}
}

func (s *TransactionService) TopUp(ctx context.Context, userID string, amount int64) (*model.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("find user: %w", err)
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore + amount

	tx := &model.Transaction{
		TransactionID:   uuid.New().String(),
		UserID:          userID,
		TransactionType: model.TypeCredit,
		TransactionKind: model.KindTopup,
		Amount:          amount,
		Status:          model.StatusSuccess,
		BalanceBefore:   balanceBefore,
		BalanceAfter:    balanceAfter,
	}

	if err := s.db.Transaction(func(txDB *gorm.DB) error {
		if err := txDB.Create(tx).Error; err != nil {
			return err
		}
		if err := txDB.Model(&model.User{}).
			Where("user_id = ?", userID).
			Update("balance", balanceAfter).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("topup transaction: %w", err)
	}

	return tx, nil
}

func (s *TransactionService) Payment(ctx context.Context, userID string, amount int64, remarks string) (*model.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	tx := &model.Transaction{}
	err := s.db.Transaction(func(txDB *gorm.DB) error {
		var user model.User
		if err := txDB.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", userID).
			First(&user).Error; err != nil {
			return err
		}

		if user.Balance < amount {
			return ErrInsufficientBalance
		}

		balanceBefore := user.Balance
		balanceAfter := balanceBefore - amount

		tx.TransactionID = uuid.New().String()
		tx.UserID = userID
		tx.TransactionType = model.TypeDebit
		tx.TransactionKind = model.KindPayment
		tx.Amount = amount
		tx.Remarks = remarks
		tx.Status = model.StatusSuccess
		tx.BalanceBefore = balanceBefore
		tx.BalanceAfter = balanceAfter

		if err := txDB.Create(tx).Error; err != nil {
			return err
		}

		if err := txDB.Model(&model.User{}).
			Where("user_id = ?", userID).
			Update("balance", balanceAfter).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("payment transaction: %w", err)
	}

	return tx, nil
}

func (s *TransactionService) Transfer(ctx context.Context, senderID, receiverID string, amount int64, remarks string) (*model.Transaction, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	_, err := s.userRepo.FindByID(ctx, receiverID)
	if err != nil {
		return nil, fmt.Errorf("find receiver: %w", err)
	}

	tx := &model.Transaction{
		TransactionID:   uuid.New().String(),
		UserID:          senderID,
		TransactionType: model.TypeDebit,
		TransactionKind: model.KindTransfer,
		Amount:          amount,
		Remarks:         remarks,
		Status:          model.StatusPending,
		RelatedUserID:   receiverID,
	}

	if err := s.txRepo.Create(ctx, tx); err != nil {
		return nil, fmt.Errorf("create transfer transaction: %w", err)
	}

	s.transferQ <- worker.TransferJob{
		TransactionID: tx.TransactionID,
		SenderID:     senderID,
		ReceiverID:   receiverID,
		Amount:       amount,
	}

	return tx, nil
}

func (s *TransactionService) GetTransactions(ctx context.Context, userID string) ([]model.Transaction, error) {
	txns, err := s.txRepo.FindByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get transactions: %w", err)
	}
	return txns, nil
}
