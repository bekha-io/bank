package handlers

import (
	"banking/internal/http/reqModels"
	"banking/pkg/errors"
	"banking/pkg/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

var sessionUserKey = "contextUserId"

func CreateUser(c *gin.Context) {
	var resp = JSONResp{
		Status: RespStatusFail,
	}

	// Грязный input
	var req reqModels.CreateUser
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Error = err.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	var user *models.User
	var err error

	// Here we are creating user depending on his status
	if req.LegalEntityInfo != nil {

		user, err = service.CreateUserForLegalEntity(req.Login, req.Password, req.PhoneNumber, req.LegalEntityInfo.OfficialName,
			req.LegalEntityInfo.TaxIdNumber)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
	} else if req.IndividualInfo != nil {
		user, err = service.CreateUser(req.Login, req.Password, req.PhoneNumber, req.IndividualInfo.Name,
			req.IndividualInfo.LastName, req.IndividualInfo.MiddleName)
		if err != nil {
			resp.Error = err.Error()
			c.JSON(http.StatusBadRequest, resp)
			return
		}
	} else {
		resp.Error = errors.ShouldBeIndividualOrLegalEntity.Error()
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	resp.Status = RespStatusOK
	resp.Body = user
	c.JSON(http.StatusCreated, resp)
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

//func LoginPage(c *gin.Context) {
//	session := sessions.Default(c)
//
//	userId := session.Get(sessionUserKey)
//	if userId != nil {
//		Logout(c)
//	}
//
//}

//func Logout(c *gin.Context) {
//	session := sessions.Default(c)
//}

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	//session := sessions.Default(c)
	//user := session.Get(sessionUserKey)
	//if user == nil {
	//	// Abort the request with the appropriate error code
	//	c.Redirect(http.StatusUnauthorized, "/login")
	//	return
	//}
	// Continue down the chain to handler etc
	c.Next()
}

func setupAuthHandlers(r *gin.Engine) {
	r.POST("/createUser", CreateUser)
	r.POST("/refreshToken", RefreshToken)

	//cookieSecret := []byte("billieeilish")
	//client := r.Group("")
	//client.Use(sessions.Sessions("everythingiwanted", sessions.NewCookieStore(cookieSecret)))
	//client.GET("/login", LoginPage)
}
