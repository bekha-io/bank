package models

import (
	"banking/pkg/types"
	"errors"
	gonanoid "github.com/matoous/go-nanoid"
	"gorm.io/gorm"
	"time"
)

type Payment struct {
	BaseModel
	ID                string `gorm:"primaryKey" json:"id"`
	MerchantUserID    string `gorm:"notNull" json:"merchant_user_id"`
	MerchantAccountID string `gorm:"notNull" json:"account_id"`
	MerchantPaymentID string `json:"merchant_payment_id,omitempty"`

	Status    types.TransactionStatus `json:"status"`
	ExpiresAt *time.Time              `json:"expires_at"`

	Amount   types.Money    `gorm:"notNull" json:"amount"`
	Currency types.Currency `gorm:"notNull" json:"currency"`
	Purpose  string         `json:"payment_purpose,omitempty"`

	IsMultiple bool `json:"is_multiple"`

	Error string `json:"error,omitempty"`

	SuccessUrl string `json:"success_url,omitempty"`
	FailUrl    string `json:"fail_url,omitempty"`
}

func (p Payment) isValidOnCreate() (err error) {
	if p.Amount < 100 {
		return errors.New("payment amount should be greater than 100")
	}

	// Whether the merchant is a legal entity
	var user User
	Db.First(&user, "id = ?", p.MerchantUserID)
	if user.LegalStatus != types.LegalEntityStatus {
		return errors.New("individuals cannot issue payments")
	}

	// Whether the provided account id is owned by the user
	res := Db.Where(&Account{ID: p.MerchantAccountID, UserID: p.MerchantUserID}).Find(&Account{})
	if res.RowsAffected == 0 {
		return errors.New("an account that are bind to payment should be owned by the user that has generated a payment")
	}

	return err
}

func (p *Payment) BeforeCreate(tx *gorm.DB) (err error) {
	err = p.isValidOnCreate()

	// If expires at was not defined
	if p.ExpiresAt == nil {
		var expiresAt = time.Now().AddDate(0, 0, 3)
		p.ExpiresAt = &expiresAt
	}

	p.ID = gonanoid.MustID(54)
	p.Status = types.TransactionStatusPending
	return
}

func (p *Payment) IsActive() bool {
	return time.Now().Before(*p.ExpiresAt) && p.Status != types.TransactionStatusProcessed
}

func NewPayment(merchantUserId, merchantAccountId string, amount types.Money, currency types.Currency, isMultiple bool) *Payment {
	return &Payment{
		MerchantUserID:    merchantUserId,
		MerchantAccountID: merchantAccountId,
		Status:            types.TransactionStatusDraft,
		Amount:            amount,
		Currency:          currency,
		IsMultiple:        isMultiple,
	}
}

func (p *Payment) SetMultiple() {
	p.IsMultiple = true
}
