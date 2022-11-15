package types

import (
	"banking/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
	"time"
	"unicode/utf8"
)

type Currency string

func (c Currency) String() string {
	return string(c)
}

var Digits = "0123456789"

const (
	TJS Currency = "TJS"
	USD Currency = "USD"
	RUB Currency = "RUB"
	EUR Currency = "EUR"
	TON Currency = "TON"
)

func (c Currency) MinTransferAmount() Money {
	switch c {
	case TJS:
		return Money(100 * 1) // 1 TJS
	case USD:
		return Money(100 * 1) // 1 USD
	case RUB:
		return Money(100 * 1) // 1 RUB
	case EUR:
		return Money(100 * 1) // 1 EUR
	case TON:
		return Money(50 * 1) // 0.5 TON
	default:
		return Money(100 * 1) // 1 major in any other currency
	}
}

var AllowedCurrencies = []string{TJS.String(), USD.String(), RUB.String(), EUR.String()}

type Money int64

type PAN string
type PIN string
type CV2 string
type ExpireDate string

type CardSystem string

const (
	Visa       CardSystem = "VISA"
	MasterCard CardSystem = "MC"
	KortiMilli CardSystem = "KM"
)

type TransactionStatus string

const (
	TransactionStatusDraft     TransactionStatus = "Draft"
	TransactionStatusCanceled  TransactionStatus = "Canceled"  // By customer
	TransactionStatusPending   TransactionStatus = "Pending"   // Waiting
	TransactionStatusProcessed TransactionStatus = "Processed" // OK
	TransactionStatusRefused   TransactionStatus = "Refused"   // By us
	//TransactionStatusConfirmation TransactionStatus = "Confirmation" // Waiting for confirmation via 3D-Secure
)

type TransactionType int

const (
	Debit  TransactionType = 1
	Credit TransactionType = 0
)

type PhoneNumber string

func (ph PhoneNumber) Validate() (err error) {
	// Checking phone number is valid (starts with 992 and contains 12 digits)
	var phoneCountryCode = ph[:3]
	if phoneCountryCode != "992" || utf8.RuneCountInString(string(ph)) != 12 {
		return errors.InvalidPhoneNumberError
	}

	return nil
}

func (p PAN) Validate() (err error) {
	if utf8.RuneCountInString(string(p)) != 16 {
		return errors.ValidatePANError
	}
	return
}

func (p PIN) Validate() (err error) {
	if utf8.RuneCountInString(string(p)) != 3 {
		return errors.ValidatePINError
	}
	return
}

func (c CV2) Validate() (err error) {
	if utf8.RuneCountInString(string(c)) != 3 {
		return errors.ValidateCV2Error
	}
	return
}

func (e ExpireDate) AsTime() (*time.Time, error) {
	tm, err := time.Parse("01/06", string(e))
	if err != nil {
		return nil, err
	}
	return &tm, err
}

func (e ExpireDate) Validate() (err error) {
	if _, err := e.AsTime(); err != nil {
		return err
	}
	return
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
