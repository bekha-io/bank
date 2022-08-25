package reqModels

import "banking/pkg/types"

type AllUserCards struct {
	Login string `json:"login" binding:"required"`
}

type IssueCard struct {
	AccountID  string           `json:"accountId" binding:"required"`
	CardSystem types.CardSystem `json:"cardSystem"`
}
