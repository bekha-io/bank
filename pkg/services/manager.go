package services

import (
	"banking/pkg/models"
	"gorm.io/gorm"
)

type ServiceManager struct {
	Card        cardService
	User        userService
	Account     accountService
	Transaction transactionService
	Auth        authInterface

	db *gorm.DB
}

func NewServiceManager() *ServiceManager {
	return &ServiceManager{
		db: models.Db,
	}
}
