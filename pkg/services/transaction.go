package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
)

type transactionService interface {
	makeTransactionRecord(accountID string, amount types.Money, transactionType types.TransactionType,
		comment string) (err error)
	makeCreditRecord(accountID string, amount types.Money, comment string) (err error)
	makeDebitRecord(accountID string, amount types.Money, comment string) (err error)
}

func (s *ServiceManager) makeTransactionRecord(accountID string, amount types.Money,
	transactionType types.TransactionType,
	comment string) (trn *models.Transaction, err error) {

	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	if account != nil {
		trn = models.NewTransaction(account.ID, amount, account.Currency, transactionType)
		s.DB.Save(&trn)
		return trn, nil
	}

	return
}

func (s *ServiceManager) makeDebitRecord(accountID string, amount types.Money,
	comment string) (trn *models.Transaction, err error) {
	trn, err = s.makeTransactionRecord(accountID, amount, types.Debit, comment)
	return
}

func (s *ServiceManager) makeCreditRecord(accountID string, amount types.Money,
	comment string) (trn *models.Transaction, err error) {
	trn, err = s.makeTransactionRecord(accountID, amount, types.Credit, comment)
	return
}

func (s *ServiceManager) getTransactionByID(trnId uint) (trn *models.Transaction, err error) {
	tx := s.DB.First(&trn, "id = ?", trnId)
	if tx.RowsAffected == 0 || tx.Error != nil {
		return nil, errors.TransactionNotFound
	}
	return trn, err
}

func (s *ServiceManager) TransferMoney(origAccountId, destAccountId string, amount types.Money, comment string) (trn *models.Transfer, err error) {
	origAccount, err := s.GetAccountByID(origAccountId)
	if err != nil {
		return nil, err
	}

	destAccount, err := s.GetAccountByID(destAccountId)
	if err != nil {
		return nil, err
	}

	if origAccount.Currency != destAccount.Currency {
		return nil, errors.CurrenciesShouldMatch
	}

	if origAccount.Balance < amount {
		return nil, errors.InsufficientBalance
	}

	trn = models.NewTransfer(origAccount.ID, destAccount.ID, amount, origAccount.Currency, comment)
	res := s.DB.Create(&trn)
	if res.RowsAffected == 0 || res.Error != nil {
		return nil, res.Error
	}

	return trn, nil
}
