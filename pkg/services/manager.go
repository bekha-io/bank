package services

import (
	"banking/pkg/models"
	"gorm.io/gorm"
)

type ServiceManager struct {
	cardService
	userService
	accountService
	transactionService
	authInterface
	paymentInterface

	DB *gorm.DB
}

func NewServiceManager() *ServiceManager {
	db := models.ConnectDB()

	return &ServiceManager{
		DB: db,
	}
}
