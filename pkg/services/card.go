package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"log"
)

type cardService interface {
	IssueCard(accountID string, cardSystem types.CardSystem) (card *models.Card, err error)
}

func (s *ServiceManager) getCardByPan(pan types.PAN) (card *models.Card, err error) {
	var c models.Card
	res := s.DB.First(&c, "pan = ?", string(pan))
	if res.RowsAffected == 0 || res.Error != nil {
		return nil, res.Error
	}
	return &c, err
}

func (s *ServiceManager) IssueCard(accountID string, cardSystem types.CardSystem) (card *models.Card, err error) {
	account, err := s.GetAccountByID(accountID)
	if err != nil {
		return nil, err
	}

	card = &models.Card{
		AccountID:  account.ID,
		CardSystem: cardSystem,
	}

	// Limit cards for each account (1 per account)
	if len(account.GetCards()) >= 1 {
		return nil, errors.CardsPerAccountLimitError
	}

	if err = s.DB.Create(card).Error; err != nil {
		log.Println("Cannot insert card to DB! Error: ", err)
		return nil, err
	}

	return card, nil
}
