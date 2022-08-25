package reqModels

import "banking/pkg/types"

type AllUserAccounts struct {
	Login string `json:"login" binding:"required"`
}

type IssueAccount struct {
	Login    string         `json:"login" binding:"required"`
	Currency types.Currency `json:"currency" binding:"required"`
}
