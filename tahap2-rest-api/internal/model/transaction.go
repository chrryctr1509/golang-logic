package model

import (
	"time"
)

const (
	TypeCredit  = "CREDIT"
	TypeDebit   = "DEBIT"
)

const (
	KindTopup   = "TOPUP"
	KindPayment = "PAYMENT"
	KindTransfer = "TRANSFER"
)

const (
	StatusPending = "PENDING"
	StatusSuccess = "SUCCESS"
	StatusFailed  = "FAILED"
)

type Transaction struct {
	TransactionID   string    `gorm:"column:transaction_id;type:char(36);primaryKey"`
	UserID          string    `gorm:"column:user_id;type:char(36)"`
	TransactionType string    `gorm:"column:transaction_type;type:varchar(10)"`
	TransactionKind string    `gorm:"column:transaction_kind;type:varchar(20)"`
	Amount          int64     `gorm:"column:amount;type:bigint"`
	Remarks         string    `gorm:"column:remarks;type:text"`
	BalanceBefore   int64     `gorm:"column:balance_before;type:bigint"`
	BalanceAfter    int64     `gorm:"column:balance_after;type:bigint"`
	Status          string    `gorm:"column:status;type:varchar(20);default:SUCCESS"`
	RelatedUserID   string    `gorm:"column:related_user_id;type:char(36)"`
	CreatedDate     time.Time `gorm:"column:created_date;autoCreateTime"`
}

func (Transaction) TableName() string {
	return "transactions"
}
