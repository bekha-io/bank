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

	db *gorm.DB
}

func NewServiceManager() *ServiceManager {
	db := models.ConnectDB()

	return &ServiceManager{
		db: db,
	}
}
