package handlers

import (
	"banking/internal/http/reqModels"
	"banking/pkg/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func AllUserCards(c *gin.Context) {
	var req reqModels.AllUserCards
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

	cards := user.GetCards()
	for i, _ := range cards {
		cards[i].Mask()
	}
	resp.Status = RespStatusOK
	resp.Body = cards
	c.IndentedJSON(http.StatusOK, resp)
}

func IssueCard(c *gin.Context) {
	var req reqModels.IssueCard
	var resp = DefaultResp

	if err := c.BindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	account, err := service.GetAccountByID(req.AccountID)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusNotFound, resp)
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
	user, err := service.GetUserById(account.UserID)
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

	issuedCard, err := service.IssueCard(account.ID, req.CardSystem)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	resp.Body = issuedCard
	resp.Status = RespStatusOK
	c.IndentedJSON(http.StatusCreated, resp)
	return
}

func setupCardsHandlers(r *gin.Engine) {
	g := r.Group("")
	g.Use(TokenAuthMiddleware())

	g.GET("cards", AllUserCards)
	g.POST("cards/issue", IssueCard)
}
