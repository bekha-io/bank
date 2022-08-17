package models

import (
	"banking/pkg/types"
	"errors"
	"gorm.io/gorm"
	"math/rand"
	"time"
)

type Account struct {
	BaseModel
	ID       string         `gorm:"primaryKey;<-:create" json:"id"`
	UserID   string         `gorm:"<-:create" json:"ownerId"`
	Currency types.Currency `gorm:"default:TJS" json:"currency"`
	Balance  types.Money    `gorm:"default:0" json:"balance"`

	IsFrozen bool `gorm:"default:false;" json:"isFrozen"`
}

// Generates a random ID
func generateAccountID() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, 20)
	for i := range b {
		b[i] = types.Digits[rand.Intn(len(types.Digits))]
	}
	return string(b)
}

func (a *Account) isValidOnCreate() (err error) {
	if a.UserID == "" {
		err = errors.New("cannot create an account without owner")
	}
	return
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	a.ID = generateAccountID()
	err = a.isValidOnCreate()
	return
}

func (a *Account) GetCards() []Card {
	var cards []Card
	Db.Find(&cards, "account_id = ?", a.ID)
	return cards
}

func (a *Account) GetTransactions() []Transaction {
	var transactions []Transaction
	Db.Find(&transactions, "account_id = ?", a.ID)
	return transactions
}
