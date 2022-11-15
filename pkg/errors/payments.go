package errors

import "errors"

var (
	PaymentShouldBeDraft              = errors.New("payment status should be draft to create payment")
	IndividualsCannotGeneratePayments = errors.New("individuals cannot generate payments. Only legal entities are allowed to")
	PaymentNotFound                   = errors.New("payment with such attributes not found")
	PaymentExpired                    = errors.New("payment expired")
)
