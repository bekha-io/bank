package handlers

import (
	"banking/internal/http/reqModels"
	"banking/pkg/errors"
	"banking/pkg/models"
	"banking/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AllUserAccounts(c *gin.Context) {
	var req reqModels.AllUserAccounts
	var resp = DefaultResp

	if err := c.BindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	ctxLogin, exists := c.Get("ctxLogin")
	if !exists {
		resp.Error = errors.InvalidAccessTokenError.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	// Context user (one that are making request)
	ctxUser, err := service.GetUserByLogin(ctxLogin.(string))
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// User that are being requested
	user, err := service.GetUserByLogin(req.Login)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// Только админы и саппорт могут читать информацию других пользователей
	if ctxUser.Login != user.Login {
		if ctxUser.IsUserRole() {
			resp.Error = errors.NotEnoughRights.Error()
			c.IndentedJSON(http.StatusForbidden, resp)
			return
		}
	}

	accounts := user.GetAccounts()
	resp.Status = RespStatusOK
	resp.Body = accounts
	c.IndentedJSON(http.StatusOK, resp)
}

func TransferMoney(c *gin.Context) {
	var resp = DefaultResp
	var req reqModels.TransferMoneyReq

	ctxLogin, _ := c.Get("ctxLogin")
	user, err := service.GetUserByLogin(ctxLogin.(string))

	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	origAccount, err := service.GetAccountByID(req.OrigAccountId)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Checking does OrigAccount belong to user
	if origAccount.UserID != user.ID && user.Role != types.AdminRole {
		resp.Error = "Forbidden"
		c.JSON(http.StatusForbidden, resp)
		return
	}

	var destAccount *models.Account
	var destUser *models.User

	// If destination requisite is a AccountID
	if req.DestAccountId != nil {
		destAccount, err = service.GetAccountByID(*req.DestAccountId)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		// Else if destination requisite is phone number
	} else if req.DestPhoneNumber != nil {
		// Validating phone number
		if err = req.DestPhoneNumber.Validate(); err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		destUser, err = service.GetUserByPhoneNumber(*req.DestPhoneNumber)
		// If there is no account with such phone number
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}

		destUserAccounts := destUser.GetAccounts()
		for _, acc := range destUserAccounts {
			if acc.Currency == origAccount.Currency {
				destAccount = &acc
				break
			}
		}
		// If destUser has no accounts with the same currency as origAccount
		if destAccount == nil {
			resp.Error = "Destination account is not the same currency as origin account"
			c.JSON(http.StatusBadRequest, resp)
			return
		}
	}

	trn, err := service.TransferMoney(origAccount.ID, destAccount.ID, req.Amount, req.Comment)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Body = trn
	resp.Status = RespStatusOK
	c.JSON(http.StatusCreated, resp)
	return
}

func UserAccount(c *gin.Context) {
	var resp = DefaultResp

	accountId := c.Param("id")

	ctxLogin, exists := c.Get("ctxLogin")
	if !exists {
		resp.Error = errors.InvalidAccessTokenError.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	// Context user (one that are making request)
	ctxUser, err := service.GetUserByLogin(ctxLogin.(string))
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// Account that are being requested
	account, err := service.GetAccountByID(accountId)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// Проверка прав доступа
	if account.UserID != ctxUser.ID {
		if ctxUser.IsUserRole() {
			resp.Error = errors.NotEnoughRights.Error()
			c.IndentedJSON(http.StatusForbidden, resp)
			return
		}
	}

	resp.Status = RespStatusOK
	resp.Body = account
	c.IndentedJSON(http.StatusOK, resp)
}

func FreezeAccount(c *gin.Context) {

}

func UnfreezeAccount(c *gin.Context) {

}

func IssueAccount(c *gin.Context) {
	var req reqModels.IssueAccount
	var resp = DefaultResp

	if err := c.BindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	ctxLogin, exists := c.Get("ctxLogin")
	if !exists {
		resp.Error = errors.InvalidAccessTokenError.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	// Context user (one that are making request)
	ctxUser, err := service.GetUserByLogin(ctxLogin.(string))
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// User that are being requested
	user, err := service.GetUserByLogin(req.Login)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
		return
	}

	// Проверка прав доступа
	if user.ID != ctxUser.ID {
		if ctxUser.IsUserRole() {
			resp.Error = errors.NotEnoughRights.Error()
			c.IndentedJSON(http.StatusForbidden, resp)
			return
		}
	}

	newAccount, err := service.IssueAccount(user.ID, req.Currency)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	resp.Body = newAccount
	resp.Status = RespStatusOK
	c.IndentedJSON(http.StatusCreated, resp)
	return
}

func setupAccountsHandlers(r *gin.Engine) {
	g := r.Group("/")
	g.Use(TokenAuthMiddleware())
	g.GET("accounts", AllUserAccounts)
	g.GET("accounts/:id", UserAccount)
	g.POST("accounts/issue", IssueAccount)
	g.POST("accounts/:id/freeze", FreezeAccount)
	g.POST("accounts/:id/unfreeze", UnfreezeAccount)
	g.POST("accounts/transfer", TransferMoney)
}
