package services

import (
	"banking/pkg/models"
	"banking/pkg/types"
)

type accountService interface {
	IssueAccount(userID string, currency types.Currency) (account *models.Account, err error)
	GetAccountByID(accountID string) (account *models.Account)
	FreezeAccount(accountID string) (ok bool)
	UnfreezeAccount(accountID string) (ok bool)
}

func (s *ServiceManager) IssueAccount(userID string, currency types.Currency) (account *models.Account, err error) {
	account = &models.Account{
		UserID:   userID,
		Currency: currency,
		Balance:  0,
	}
	result := s.db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}
	return account, nil
}

func (s *ServiceManager) GetAccountByID(accountID string) (account *models.Account) {
	s.db.First(&account, "id = ?", accountID)

	if account == (&models.Account{}) {
		return nil
	}

	return account
}

func (s *ServiceManager) FreezeAccount(accountID string) (ok bool) {
	account := s.GetAccountByID(accountID)

	if account == nil {
		return false
	}

	account.IsFrozen = true
	s.db.Save(account)
	return true
}

func (s *ServiceManager) UnfreezeAccount(accountID string) (ok bool) {
	account := s.GetAccountByID(accountID)

	if account == nil {
		return false
	}
	account.IsFrozen = false
	s.db.Save(account)
	return true
}
