package types

import (
	"banking/pkg/errors"
	"github.com/golang-jwt/jwt/v4"
)

type Currency string

var Digits = "0123456789"

const (
	TJS Currency = "TJS"
	USD Currency = "USD"
	RUB Currency = "RUB"
)

type Money int64

type PAN string

type CardSystem string

const (
	Visa       CardSystem = "VISA"
	MasterCard CardSystem = "MC"
	KortiMilli CardSystem = "KM"
)

type TransactionStatus string

const (
	TransactionStatusCanceled  TransactionStatus = "Canceled"  // By customer
	TransactionStatusPending   TransactionStatus = "Pending"   // Waiting
	TransactionStatusProcessed TransactionStatus = "Processed" // OK
	//TransactionStatusRefused      TransactionStatus = "Refused"      // By us
	//TransactionStatusConfirmation TransactionStatus = "Confirmation" // Waiting for confirmation via 3D-Secure
)

type TransactionType int

const (
	Debit  TransactionType = 1
	Credit TransactionType = 0
)

type PhoneNumber string

func (ph PhoneNumber) IsValid() (err error) {
	// Checking phone number is valid (starts with 992 and contains 12 digits)
	var phoneCountryCode = ph[:3]
	if phoneCountryCode != "992" || len(phoneCountryCode) != 12 {
		err = errors.InvalidPhoneNumberError
		return
	}

	return nil
}

type Claims struct {
	Login string `json:"login"`
	jwt.RegisteredClaims
}
