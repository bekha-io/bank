package respModels

import "banking/pkg/models"

type AccountTransactions struct {
	models.Account
	Transactions []models.Transaction `json:"transactions"`
}
