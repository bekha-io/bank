package services

import (
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"fmt"
)

type paymentInterface interface {
	GeneratePayment(paymentDraft models.Payment) (payment *models.Payment, err error)
	GetPaymentByID(paymentId string) (payment *models.Payment, err error)
}

func (s *ServiceManager) GetPaymentByID(paymentId string) (payment *models.Payment, err error) {
	res := s.DB.First(&payment, "id = ?", paymentId)
	if res.RowsAffected == 0 || res.Error != nil {
		return nil, errors.PaymentNotFound
	}
	return payment, err
}

func (s *ServiceManager) GeneratePayment(paymentDraft models.Payment) (payment *models.Payment, err error) {
	if paymentDraft.Status != types.TransactionStatusDraft {
		return nil, errors.PaymentShouldBeDraft
	}

	// Getting merchantUser
	merchantUser, err := s.GetUserById(paymentDraft.MerchantUserID)
	if err != nil {
		return nil, err
	}

	// Getting merchantsAccount
	merchantAccount, err := s.GetAccountByID(paymentDraft.MerchantAccountID)
	if err != nil {
		return nil, err
	}

	if merchantUser.LegalStatus != types.LegalEntityStatus {
		return nil, errors.IndividualsCannotGeneratePayments
	}

	// Only account's owner has permission to generate payments
	if merchantAccount.UserID != merchantUser.ID {
		return nil, errors.NotEnoughRights
	}

	// Currency cannot differ from merchant account's currency
	paymentDraft.Currency = merchantAccount.Currency

	payment = &paymentDraft
	// If everything is ok
	res := s.DB.Create(payment)
	if res.Error != nil || res.RowsAffected == 0 {
		return nil, res.Error
	}

	return payment, nil
}

func (s *ServiceManager) Pay(paymentId string,
	pan types.PAN, expireDate types.ExpireDate, cv2 types.CV2) (trn *models.Transaction, err error) {

	if err := pan.Validate(); err != nil {
		return nil, err
	}

	if err = cv2.Validate(); err != nil {
		return nil, err
	}

	if err = expireDate.Validate(); err != nil {
		return nil, err
	}

	card, err := s.getCardByPan(pan)
	if err != nil {
		return nil, errors.CardNotFound
	}

	if !card.CompareWith(pan, expireDate, cv2) {
		return nil, errors.CardNotFound
	}

	payment, err := s.GetPaymentByID(paymentId)
	if err != nil {
		return nil, errors.PaymentNotFound
	}

	if !payment.IsActive() {
		return nil, errors.PaymentExpired
	}

	trf, err := s.TransferMoney(card.AccountID, payment.MerchantAccountID, payment.Amount,
		fmt.Sprintf("Оплата мерчанту %v", payment.MerchantAccountID))
	if err != nil {
		return nil, err
	}

	trf.ReferencePaymentID = &payment.ID
	s.DB.Save(&trf)

	trn, err = s.getTransactionByID(*trf.OrigTransactionID)
	if err != nil {
		return nil, err
	}

	return trn, nil
}
