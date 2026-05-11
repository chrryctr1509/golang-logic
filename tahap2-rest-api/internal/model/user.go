package model

import (
	"time"
)

type User struct {
	UserID      string    `gorm:"column:user_id;type:char(36);primaryKey"`
	FirstName   string    `gorm:"column:first_name;type:varchar(100)"`
	LastName    string    `gorm:"column:last_name;type:varchar(100)"`
	PhoneNumber string    `gorm:"column:phone_number;type:varchar(20);uniqueIndex"`
	Address     string    `gorm:"column:address;type:text"`
	PIN         string    `gorm:"column:pin;type:varchar(255)"`
	Balance     int64     `gorm:"column:balance;type:bigint;default:0"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
	UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}

func (User) TableName() string {
	return "users"
}
