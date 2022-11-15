package jobs

import (
	"banking/pkg/models"
	"banking/pkg/services"
	"banking/pkg/types"
	"fmt"
	"log"
	"time"
)

var svc *services.ServiceManager

func TransactionValidationJob() {
	log.Println("Running TRANSACTION (non-transfer) job ...")

	svc = services.NewServiceManager()
	for {
		// Ranging over non-transfer transactions only
		var transactions []models.Transaction
		svc.DB.Model(&models.Transaction{}).Where("status = ?", types.TransactionStatusPending).Where(" reference_trn_id <> null").Scan(&transactions)

		for _, trn := range transactions {
			log.Printf("[TRN %v] Checking transaction (%v %v) ...", trn.ID, trn.Amount, trn.Currency)

			account, err := svc.GetAccountByID(trn.AccountID)
			if err != nil {
				trn.Error = "Invalid account ID"
				trn.Status = types.TransactionStatusRefused
				svc.DB.Save(&trn)
				continue
			}

			// Currencies do not match
			if account.Currency != trn.Currency {
				trn.Status = types.TransactionStatusRefused
				trn.Error = "Transaction's currency does not match account's currency"
				svc.DB.Save(&trn)
				continue
			}

			// Insufficient balance if transaction's type is Credit
			if account.Balance < trn.Amount && trn.Type == types.Credit {
				trn.Status = types.TransactionStatusRefused
				trn.Error = "Insufficient balance"
				svc.DB.Save(&trn)
				continue
			}

			switch trn.Type {
			case types.Debit:
				account.Balance += trn.Amount
			case types.Credit:
				account.Balance -= trn.Amount
			}

			trn.Status = types.TransactionStatusProcessed
			svc.DB.Save(&trn)
			svc.DB.Save(&account)
		}

		time.Sleep(3)
	}
}

func nullifyTransfer(transfer *models.Transfer, error string, status types.TransactionStatus) {
	transfer.Status = status
	transfer.Error = error
	svc.DB.Save(&transfer)
	svc.DB.Model(&models.Transaction{}).Where("id IN ?",
		[]uint{*transfer.OrigTransactionID, *transfer.DestTransactionID}).Updates(models.Transaction{Status: status})
}

func markSuccessTransfer(transfer *models.Transfer) {
	transfer.Status = types.TransactionStatusProcessed
	svc.DB.Save(&transfer)
	svc.DB.Model(&models.Transaction{}).Where("id IN ?",
		[]uint{*transfer.OrigTransactionID, *transfer.DestTransactionID}).Updates(models.Transaction{Status: types.TransactionStatusProcessed})
}

func TransferValidationJob() {
	log.Println("Running TRANSFER job...")

	svc = services.NewServiceManager()
	for {
		var transfers []models.Transfer
		svc.DB.Find(&transfers, "status = ?", types.TransactionStatusPending)

		for _, trf := range transfers {
			log.Printf("[TRF %v] validating", trf.ID)

			// Getting related transactions
			var origTrn models.Transaction
			var destTrn models.Transaction

			var condMatch bool = true

			res := svc.DB.Find(&origTrn, "id = ?", trf.OrigTransactionID)
			if res.RowsAffected == 0 {
				condMatch = false
			}

			res = svc.DB.Find(&destTrn, "id = ?", trf.DestTransactionID)
			if res.RowsAffected == 0 {
				condMatch = false
			}

			// If any of transactions does not exist
			if !condMatch {
				nullifyTransfer(&trf, "Origin or destination transactions do not exist", types.TransactionStatusRefused)
				continue
			}

			// If related transactions have different statuses
			if origTrn.Status != types.TransactionStatusPending || destTrn.Status != types.TransactionStatusPending {
				nullifyTransfer(&trf, fmt.Sprintf("One of related transaction has different status. (ORIG: %v, DEST: %v)", origTrn.ID, destTrn.ID),
					types.TransactionStatusRefused)
				continue
			}

			// If origin account does not exist
			origAccount, err := svc.GetAccountByID(origTrn.AccountID)
			if err != nil {
				nullifyTransfer(&trf, "Origin account with such ID does not exist", types.TransactionStatusRefused)
				continue
			}

			// If destination account does not exist
			destAccount, err := svc.GetAccountByID(destTrn.AccountID)
			if err != nil {
				nullifyTransfer(&trf, "Destination account with such ID does not exist", types.TransactionStatusRefused)
				continue
			}

			if origAccount.Currency != destAccount.Currency {
				nullifyTransfer(&trf, "Accounts' currencies do not match", types.TransactionStatusRefused)
				continue
			}

			if origAccount.Balance < trf.Amount {
				nullifyTransfer(&trf, "Insufficient balance", types.TransactionStatusRefused)
				continue
			}

			origAccount.Balance = origAccount.Balance - trf.Amount
			destAccount.Balance = destAccount.Balance + trf.Amount
			svc.DB.Save(&origAccount)
			svc.DB.Save(&destAccount)

			markSuccessTransfer(&trf)

		}
		time.Sleep(3)
	}
}

func markProcessedPayment(payment *models.Payment) {
	payment.Status = types.TransactionStatusProcessed
	svc.DB.Save(&payment)
}

func PaymentValidationJob() {
	log.Println("Running PAYMENT job ...")

	svc := services.NewServiceManager()

	for {
		var payments []*models.Payment
		svc.DB.Find(&payments, "status = ?", types.TransactionStatusPending)

		for _, payment := range payments {

			// Marking all expired payments as processed
			if time.Now().After(*payment.ExpiresAt) {
				markProcessedPayment(payment)
				continue
			}

			// If payment is multiple, then pass it
			if payment.IsMultiple {
				continue
			}

			// Then if payment is one-time, check whether related transfer is processed
			var relatedTrf *models.Transfer
			svc.DB.Find(&relatedTrf, "reference_payment_id = ?", payment.ID)

			if relatedTrf != nil {
				if relatedTrf.Status == types.TransactionStatusProcessed {
					markProcessedPayment(payment)
					continue
				}
			}
		}

		time.Sleep(3)
	}
}
