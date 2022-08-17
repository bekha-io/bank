package handlers

import (
	"banking/internal/http/json/reqModels"
	"github.com/gin-gonic/gin"
	"net/http"
)

func MyAccounts(c *gin.Context) {
	var req reqModels.MyAccounts
	var resp = DefaultResp

	if err := c.BindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	user, err := service.GetUserByLogin(req.Login)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	accounts := user.GetAccounts()
	resp.Status = RespStatusOK
	resp.Body = accounts
	c.IndentedJSON(http.StatusOK, resp)
}

func setupAccountsHandlers(r *gin.Engine) {
	g := r.Group("/")
	g.Use(TokenAuthMiddleware())

	g.GET("accounts", MyAccounts)
}
