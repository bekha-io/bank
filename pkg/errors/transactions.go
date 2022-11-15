package errors

import "errors"

var (
	TransactionNotFound   = errors.New("transaction with such params not found")
	CurrenciesShouldMatch = errors.New("origin account and destination account currencies should match")
)
