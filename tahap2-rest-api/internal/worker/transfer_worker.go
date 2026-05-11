package worker

import (
	"context"
	"errors"
	"log"

	"github.com/user/tahap2-rest-api/internal/model"
	"github.com/user/tahap2-rest-api/internal/repository"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TransferJob struct {
	TransactionID string
	SenderID     string
	ReceiverID   string
	Amount       int64
}

type TransferWorker struct {
	Queue  chan TransferJob
	db     *gorm.DB
	txRepo repository.TransactionRepository
}

func NewTransferWorker(db *gorm.DB, txRepo repository.TransactionRepository) *TransferWorker {
	return &TransferWorker{
		Queue:  make(chan TransferJob, 100),
		db:     db,
		txRepo: txRepo,
	}
}

func (w *TransferWorker) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("transfer worker shutting down")
			return
		case job := <-w.Queue:
			w.processTransfer(job)
		}
	}
}

func (w *TransferWorker) processTransfer(job TransferJob) {
	err := w.db.Transaction(func(tx *gorm.DB) error {
		var sender model.User
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("user_id = ?", job.SenderID).
			First(&sender).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("sender not found")
			}
			return err
		}

		if sender.Balance < job.Amount {
			return errors.New("balance is not enough")
		}

		// Debit sender
		if err := tx.Model(&model.User{}).
			Where("user_id = ?", job.SenderID).
			Update("balance", sender.Balance-job.Amount).Error; err != nil {
			return err
		}

		// Credit receiver
		if err := tx.Model(&model.User{}).
			Where("user_id = ?", job.ReceiverID).
			Update("balance", gorm.Expr("balance + ?", job.Amount)).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("transfer failed tx=%s: %v", job.TransactionID, err)
		_ = w.txRepo.UpdateStatus(context.Background(), job.TransactionID, model.StatusFailed)
		return
	}

	// Update tx status to SUCCESS
	if err := w.txRepo.UpdateStatus(context.Background(), job.TransactionID, model.StatusSuccess); err != nil {
		log.Printf("failed to update tx status to SUCCESS: %v", err)
	}
	log.Printf("transfer completed tx=%s sender=%s receiver=%s amount=%d",
		job.TransactionID, job.SenderID, job.ReceiverID, job.Amount)
}