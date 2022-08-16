package handlers

import (
	"banking/pkg/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createUserReq struct {
	Login       string            `json:"login"`
	Password    string            `json:"password"`
	FirstName   string            `json:"firstName"`
	LastName    string            `json:"lastName"`
	MiddleName  string            `json:"middleName"`
	PhoneNumber types.PhoneNumber `json:"phoneNumber"`
}

type refreshTokenReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func CreateUser(c *gin.Context) {
	var resp = JSONResp{
		Status: RespStatusFail,
	}

	// Грязный input
	var reqBody createUserReq
	if err := c.BindJSON(&reqBody); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
	}

	// Пользователя не удалось создать
	user, err := service.User.CreateUser(reqBody.Login, reqBody.Password, reqBody.PhoneNumber, reqBody.FirstName,
		reqBody.LastName, reqBody.MiddleName)
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusInternalServerError, resp)
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
	var reqBody refreshTokenReq
	if err := c.BindJSON(&reqBody); err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusBadRequest, resp)
	}

	// Сервис не может выпустить новый токен
	accessToken, err := service.Auth.RefreshAccessToken(reqBody.Login, []byte(reqBody.Password))
	if err != nil {
		resp.Error = err.Error()
		c.IndentedJSON(http.StatusInternalServerError, resp)
	}

	resp.Status = RespStatusOK
	resp.Body = map[string]string{"access_token": accessToken}
	c.IndentedJSON(http.StatusOK, resp)
}

func setupAuthHandlers(r *gin.Engine) {
	r.POST("/createUser", CreateUser)
	r.POST("/refreshToken", RefreshToken)
}
