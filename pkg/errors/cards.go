package errors

import "errors"

var (
	CardsPerAccountLimitError = errors.New("cards amount per account limit exceeded")
	CardNotFound              = errors.New("card with such attributes not found")
)
