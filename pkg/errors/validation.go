package errors

import "errors"

var (
	ValidatePANError = errors.New("invalid pan format")
	ValidateCV2Error = errors.New("invalid CV2 format")
	ValidatePINError = errors.New("invalid PIN error")
)
