package jobs

import "log"

func RunJobs() {
	log.Println("Running jobs ...")
	go TransactionValidationJob()
	go TransferValidationJob()
	go PaymentValidationJob()
}
