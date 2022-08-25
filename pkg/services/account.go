package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
)

type accountService interface {
	IssueAccount(userID string, currency types.Currency) (account *models.Account, err error)
	GetAccountByID(accountID string) (account *models.Account, err error)
	FreezeAccount(accountID string) (ok bool)
	UnfreezeAccount(accountID string) (ok bool)
}

func (s *ServiceManager) IssueAccount(userID string, currency types.Currency) (account *models.Account, err error) {
	user, err := s.GetUserById(userID)
	if err != nil {
		return nil, err
	}

	// Если указанной валюты нет в списке разрешенных
	if !utils.Contains(types.AllowedCurrencies, string(currency)) {
		return nil, errors.UnknownCurrency
	}

	// 1 валюта - 1 счет
	for _, acc := range user.GetAccounts() {
		if acc.Currency == currency {
			return nil, errors.TwoAccountsWithSameCurrency
		}
	}

	account = &models.Account{
		UserID:   user.ID,
		Currency: currency,
		Balance:  0,
	}

	result := s.db.Create(account)
	if result.Error != nil {
		return nil, result.Error
	}

	return account, nil
}

func (s *ServiceManager) GetAccountByID(accountID string) (account *models.Account, err error) {
	res := s.db.First(&account, "id = ?", accountID)
	if res.Error != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.AccountNotFound
		}
		return nil, res.Error
	}
	return account, nil
}

func (s *ServiceManager) FreezeAccount(accountID string) (ok bool) {
	account, _ := s.GetAccountByID(accountID)

	if account == nil {
		return false
	}

	account.IsFrozen = true
	s.db.Save(account)
	return true
}

func (s *ServiceManager) UnfreezeAccount(accountID string) (ok bool) {
	account, _ := s.GetAccountByID(accountID)

	if account == nil {
		return false
	}
	account.IsFrozen = false
	s.db.Save(account)
	return true
}
