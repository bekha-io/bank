package main

import (
	"banking/internal/http/gin/handlers"
	"banking/pkg/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Here we initialize all the ORM models
	_ = models.ConnectDB()
	_ = models.Db.AutoMigrate(&models.User{}, &models.Account{}, &models.Card{}, &models.Transaction{})

	router := gin.Default()
	handlers.SetupHandlers(router)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
