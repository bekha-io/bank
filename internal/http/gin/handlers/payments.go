package handlers

import (
	"banking/internal/http/reqModels"
	"banking/pkg/errors"
	"banking/pkg/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GeneratePayment(c *gin.Context) {
	var req reqModels.GeneratePaymentReq
	var resp = DefaultResp

	ctxLogin, _ := c.Get("ctxLogin")

	// Binding request
	if err := c.BindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	// Time should not be past
	if req.ExpiresAt != nil {
		if time.Now().After(req.ExpiresAt.Local()) {
			resp.Error = errors.PaymentExpired.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
	}

	// Get context User
	user, err := service.GetUserByLogin(ctxLogin.(string))
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	paymentDraft := models.NewPayment(user.ID, req.MerchantAccountID, req.Amount, "", req.IsMultiple)
	paymentDraft.MerchantPaymentID = req.MerchantPaymentID
	paymentDraft.ExpiresAt = req.ExpiresAt
	paymentDraft.Purpose = req.Purpose

	payment, err := service.GeneratePayment(*paymentDraft)
	if err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = RespStatusOK
	resp.Body = payment
	c.JSON(http.StatusCreated, resp)
	return
}

func PaymentPage(c *gin.Context) {
	paymentId := c.Param("id")

	payment, err := service.GetPaymentByID(paymentId)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	// If payment already expired
	if payment.ExpiresAt != nil {
		if time.Now().After(*payment.ExpiresAt) {
			c.String(http.StatusNotFound, errors.PaymentExpired.Error())
			return
		}
	}

	merchant, err := service.GetUserById(payment.MerchantUserID)
	if err != nil {
		c.String(http.StatusNotFound, err.Error())
		return
	}

	merchantName, _ := merchant.CustomerInfo["official_name"].(string)

	c.HTML(http.StatusOK, "paymentPage.tmpl", gin.H{
		"payment":      payment,
		"merchantName": merchantName,
		"amount":       fmt.Sprintf("%.2f", float64(payment.Amount/100)),
	})
}

func Pay(c *gin.Context) {
	paymentId := c.Param("id")
	var form reqModels.PaymentForm

	var color string = "red"
	var message string = "Неизвестная ошибка, попробуйте еще раз!"

	// Validating form
	if err := c.ShouldBind(&form); err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	trn, err := service.Pay(paymentId, form.PAN, form.ExpireDate, form.CV2)
	if err != nil {
		message = err.Error()
	} else {
		message = "Оплата успешно произведена!"
		color = "green"
	}

	c.HTML(http.StatusOK, "paymentPost.tmpl", gin.H{
		"color":   color,
		"message": message,
		"trn":     trn,
	})
	return
}

func setupPaymentsHandlers(r *gin.Engine) {
	g := r.Group("payment")
	g.GET(":id", PaymentPage)
	g.POST(":id", Pay)
	g.POST("", TokenAuthMiddleware(), GeneratePayment)
}
