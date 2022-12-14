package handlers

import (
	"banking/pkg/services"
	"github.com/gin-gonic/gin"
)

type RespStatus string

var service = services.NewServiceManager()

const (
	RespStatusOK   RespStatus = "ok"
	RespStatusFail RespStatus = "fail"
)

type JSONResp struct {
	Status RespStatus  `json:"status"`
	Error  string      `json:"error,omitempty"`
	Body   interface{} `json:"body,omitempty"`
}

var DefaultResp = JSONResp{
	Status: RespStatusFail,
	Error:  "",
	Body:   "",
}

func SetupHandlers(r *gin.Engine) {
	setupAuthHandlers(r)
	setupAccountsHandlers(r)
	setupCardsHandlers(r)
	setupTransactionsHandlers(r)
	setupPaymentsHandlers(r)
}
