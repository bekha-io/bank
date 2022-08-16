package services

import (
	"banking/pkg/models"
	"banking/pkg/types"
)

type transactionService interface {
	makeTransactionRecord(accountID string, amount types.Money, transactionType types.TransactionType,
		comment string) (ok bool)
	makeCreditRecord(accountID string, amount types.Money, comment string) (ok bool)
	makeDebitRecord(accountID string, amount types.Money, comment string) (ok bool)
}

func (s *ServiceManager) makeTransactionRecord(accountID string, amount types.Money, transactionType types.TransactionType,
	comment string) (ok bool) {

	account := s.GetAccountByID(accountID)

	if account.ID != "" {
		var tr = &models.Transaction{
			AccountID: account.ID,
			Amount:    amount,
			Type:      transactionType,
			Status:    types.TransactionStatusProcessed,
			Comment:   comment,
		}
		s.db.Create(tr)
		return true
	}

	return false
}

func (s *ServiceManager) makeDebitRecord(accountID string, amount types.Money, comment string) (ok bool) {
	ok = s.makeTransactionRecord(accountID, amount, types.Debit, comment)
	return
}

func (s *ServiceManager) makeCreditRecord(accountID string, amount types.Money, comment string) (ok bool) {
	ok = s.makeTransactionRecord(accountID, amount, types.Credit, comment)
	return
}
