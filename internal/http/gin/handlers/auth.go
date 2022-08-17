package handlers

import (
	"banking/internal/http/json/reqModels"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateUser(c *gin.Context) {
	var resp = JSONResp{
		Status: RespStatusFail,
	}

	// Грязный input
	var reqBody reqModels.CreateUser
	if err := c.BindJSON(&reqBody); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	// Пользователя не удалось создать
	user, err := service.CreateUser(reqBody.Login, reqBody.Password, reqBody.PhoneNumber, reqBody.FirstName,
		reqBody.LastName, reqBody.MiddleName)
	if err != nil {
		fmt.Print("!@#!@#!@# ", err.Error())
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = RespStatusOK
	resp.Body = user
	c.IndentedJSON(http.StatusCreated, resp)
}

func RefreshToken(c *gin.Context) {
	var resp = JSONResp{
		Status: RespStatusFail,
	}

	// Грязный input
	var reqBody reqModels.RefreshToken
	if err := c.BindJSON(&reqBody); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
		return
	}

	// Сервис не может выпустить новый токен
	accessToken, err := service.RefreshAccessToken(reqBody.Login, []byte(reqBody.Password))
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusInternalServerError, resp)
		return
	}

	resp.Status = RespStatusOK
	resp.Body = map[string]string{"access_token": accessToken}
	c.IndentedJSON(http.StatusOK, resp)
}

func setupAuthHandlers(r *gin.Engine) {
	r.POST("/createUser", CreateUser)
	r.POST("/refreshToken", RefreshToken)
}
