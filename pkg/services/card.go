package services

import (
	"banking/pkg/models"
	"banking/pkg/types"
	"errors"
	"log"
)

type cardService interface {
	IssueCard(accountID string, cardSystem types.CardSystem) (card *models.Card, err error)
}

func (s *ServiceManager) IssueCard(accountID string, cardSystem types.CardSystem) (card *models.Card, err error) {
	account := s.GetAccountByID(accountID)
	if account == nil {
		return nil, errors.New("account with the given ID not found")
	}

	card = &models.Card{
		AccountID:  account.ID,
		CardSystem: cardSystem,
	}

	if err = s.db.Create(card).Error; err != nil {
		log.Println("Cannot insert card to db! Error: ", err)
		return nil, err
	}

	return card, nil
}
