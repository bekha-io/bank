package errors

import "errors"

var (
	AccountNotFound = errors.New("account with given ID not found")

	TwoAccountsWithSameCurrency = errors.New("already have an account with this currency")
	UnknownCurrency             = errors.New("currency is unknown. See list of available currencies")

	InsufficientBalance = errors.New("insufficient balance in origin account")
)
