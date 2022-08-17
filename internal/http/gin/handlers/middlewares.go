package handlers

import (
	"banking/pkg/errors"
	"banking/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type authCheckReq struct {
	Authorization string `header:"Authorization" binding:"required,startswith=Bearer"`
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req authCheckReq

		var resp = JSONResp{
			Status: RespStatusFail,
		}

		if err := c.ShouldBindHeader(&req); err != nil {
			resp.Error = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		if !strings.Contains(req.Authorization, "Bearer ") {
			resp.Error = errors.ShouldBeBearerTokenError.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		// Trimming word Bearer from Token
		tokenString := strings.Replace(req.Authorization, "Bearer ", "", 1)
		login, err := utils.ParseAccessToken(tokenString)
		if err != nil {
			resp.Error = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		user, err := service.GetUserByLogin(login)
		if err != nil {
			resp.Error = err.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}

		if user.AccessToken != tokenString {
			resp.Error = errors.InvalidAccessTokenError.Error()
			c.AbortWithStatusJSON(http.StatusBadRequest, resp)
			return
		}
	}

}
