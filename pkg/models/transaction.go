package models

import (
	"banking/pkg/types"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

type Transaction struct {
	BaseModel
	ID             uint                    `gorm:"primaryKey;" json:"id"`
	AccountID      string                  `gorm:"notNull;" json:"account_id"`
	Amount         types.Money             `gorm:"notNull;" json:"amount"`
	Currency       types.Currency          `gorm:"notNull" json:"currency"`
	Type           types.TransactionType   `gorm:"notNull;" json:"is_debit"`
	Status         types.TransactionStatus `gorm:"notNull;" json:"status"`
	ReferenceTrnId *uint                   `json:"reference_trn_id"`
	Comment        string                  `json:"comment,omitempty"`
	Error          string                  `json:"error,omitempty"`
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

func NewTransaction(accountID string, amount types.Money, currency types.Currency, trnType types.TransactionType) *Transaction {
	return &Transaction{
		AccountID: accountID,
		Amount:    amount,
		Type:      trnType,
		Status:    types.TransactionStatusDraft,
		Currency:  currency,
	}
}

func NewTransactionDebit(accountID string, amount types.Money, currency types.Currency) *Transaction {
	t := NewTransaction(accountID, amount, currency, types.Debit)
	return t
}

func NewTransactionCredit(accountID string, amount types.Money, currency types.Currency) *Transaction {
	t := NewTransaction(accountID, amount, currency, types.Credit)
	return t
}

type Transfer struct {
	BaseModel
	ID                 uint                    `gorm:"primaryKey" json:"id"`
	OrigAccountID      string                  `gorm:"notNull" json:"orig_account_id"`
	DestAccountID      string                  `gorm:"notNull" json:"dest_account_id"`
	OrigTransactionID  *uint                   `json:"orig_transaction_id"`
	DestTransactionID  *uint                   `json:"dest_transaction_id"`
	ReferencePaymentID *string                 `json:"reference_payment_id,omitempty"`
	Amount             types.Money             `gorm:"notNull" json:"amount"`
	Currency           types.Currency          `gorm:"notNull" json:"currency"`
	Status             types.TransactionStatus `gorm:"notNull" json:"status"`
	Error              string                  `json:"error,omitempty"`
	Comment            string                  `json:"comment,omitempty"`
}

func (t *Transfer) isValidOnCreate() (err error) {
	if t.OrigAccountID == t.DestAccountID {
		err = errors.New("cannot transfer money within the same account")
	}

	if t.Amount < t.Currency.MinTransferAmount() {
		err = errors.New(fmt.Sprintf("amount should be more than %v %v", t.Currency.MinTransferAmount(),
			t.Currency))
	}
	return
}

func (t *Transfer) BeforeCreate(tx *gorm.DB) (err error) {
	err = t.isValidOnCreate()
	t.Status = types.TransactionStatusPending
	return
}

func (t *Transfer) AfterCreate(tx *gorm.DB) (err error) {
	var origTrn *Transaction
	var destTrn *Transaction

	// Creating matching transactions
	if t.OrigTransactionID == nil {
		origTrn = NewTransactionCredit(t.OrigAccountID, t.Amount, t.Currency)
		tx.Create(&origTrn)
	}

	if t.DestTransactionID == nil {
		destTrn = NewTransactionDebit(t.DestAccountID, t.Amount, t.Currency)
		tx.Create(&destTrn)
	}

	origTrn.ReferenceTrnId = &destTrn.ID
	origTrn.Comment = t.Comment
	destTrn.ReferenceTrnId = &origTrn.ID
	destTrn.Comment = t.Comment

	tx.Save(&origTrn)
	tx.Save(&destTrn)

	t.OrigTransactionID = &origTrn.ID
	t.DestTransactionID = &destTrn.ID

	tx.Save(t)

	return
}

func NewTransfer(origAccountId, destAccountId string, amount types.Money, currency types.Currency, comment string) *Transfer {
	return &Transfer{
		OrigAccountID: origAccountId,
		DestAccountID: destAccountId,
		Amount:        amount,
		Currency:      currency,
		Status:        types.TransactionStatusDraft,
		Comment:       comment,
	}
}
