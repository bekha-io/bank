package handlers

import (
	"banking/internal/http/respModels"
	"banking/pkg/errors"
	"banking/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetAllTransactions(c *gin.Context) {
	var request struct {
		AccountID string `json:"account_id"`
	}
	var resp = DefaultResp

	if err := c.ShouldBindJSON(&request); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	ctxLogin, exists := c.Get("ctxLogin")
	if !exists {
		resp.Error = errors.InvalidAccessTokenError.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	user, err := service.GetUserByLogin(ctxLogin.(string))
	if err != nil {
		resp.Error = errors.UserDoesNotExist.Error()
		c.JSON(http.StatusNotFound, resp)
		return
	}

	account, err := service.GetAccountByID(request.AccountID)
	if err != nil {
		resp.Error = errors.AccountNotFound.Error()
		c.JSON(http.StatusNotFound, resp)
		return
	}

	if user.ID != account.UserID && !user.IsAdminRole() {
		resp.Error = errors.NotEnoughRights.Error()
		c.JSON(http.StatusForbidden, resp)
		return
	}

	var trns []models.Transaction
	service.DB.Find(&trns, "account_id = ?", account.ID)

	var respModel = respModels.AccountTransactions{Account: *account}
	respModel.Transactions = trns
	resp.Body = respModel
	resp.Status = RespStatusOK
	c.JSON(http.StatusOK, resp)
	return
}

func setupTransactionsHandlers(r *gin.Engine) {
	g := r.Group("/")
	g.Use(TokenAuthMiddleware())

	g.GET("transactions", GetAllTransactions)
}
