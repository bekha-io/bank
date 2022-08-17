package models

import (
	"banking/pkg/types"
	"errors"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	ID        uint                    `gorm:"primaryKey;" json:"id"`
	AccountID string                  `gorm:"notNull;"`
	Amount    types.Money             `gorm:"notNull;"`
	Type      types.TransactionType   `gorm:"notNull;"`
	Status    types.TransactionStatus `gorm:"notNull;"`
	Comment   string
}

func (t *Transaction) isValidOnCreate() (err error) {
	if t.Amount <= 0 {
		err = errors.New("amount should be positive (greater than 0)")
	}
	return
}

func (t *Transaction) BeforeCreate(tx *gorm.DB) (err error) {
	err = t.isValidOnCreate()
	t.Status = types.TransactionStatusPending

	return
}
