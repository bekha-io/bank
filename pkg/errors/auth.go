package errors

import "errors"

var (
	InvalidPhoneNumberError = errors.New("invalid phone number format")
	InvalidAccessTokenError = errors.New("invalid access token")
	InvalidPassword         = errors.New("invalid password")

	UserDoesNotExist = errors.New("user does not exist")

	LoginOccupiedError       = errors.New("provided login is occupied")
	PhoneNumberOccupiedError = errors.New("phone number is occupied")

	ShouldBeBearerTokenError = errors.New("headers must contain a bearer token")
	InvalidAccessToken       = errors.New("invalid access token")
)
