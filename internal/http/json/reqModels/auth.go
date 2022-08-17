package reqModels

import "banking/pkg/types"

type CreateUser struct {
	Login       string            `json:"login" binding:"required"`
	Password    string            `json:"password" binding:"required"`
	FirstName   string            `json:"firstName"   binding:"required"`
	LastName    string            `json:"lastName"    binding:"required"`
	MiddleName  string            `json:"middleName"`
	PhoneNumber types.PhoneNumber `json:"phoneNumber" binding:"required"`
}

type MyAccounts struct {
	Login string `json:"login" binding:"required"`
}

type RefreshToken struct {
	Login    string `json:"login" binding:"required"`
	Password string `json:"password" binding:"required"`
}
