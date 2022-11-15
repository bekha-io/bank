package main

import (
	"banking/internal/http/gin/handlers"
	"banking/jobs"
	"banking/pkg/models"
	"github.com/gin-gonic/gin"
)

func main() {
	// Here we initialize all the ORM models
	_ = models.ConnectDB()
	_ = models.Db.AutoMigrate(&models.User{}, &models.Account{}, &models.Card{},
		&models.Transaction{}, &models.Transfer{}, &models.Payment{})

	router := gin.Default()
	router.LoadHTMLGlob("templates/*.tmpl")
	handlers.SetupHandlers(router)

	go jobs.RunJobs()

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
