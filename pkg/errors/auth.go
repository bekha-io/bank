package errors

import "errors"

var (
	InvalidPhoneNumberError = errors.New("invalid phone number format")
	InvalidAccessTokenError = errors.New("invalid access token")
	InvalidPassword         = errors.New("invalid password")

	UserDoesNotExist = errors.New("user does not exist")
)
