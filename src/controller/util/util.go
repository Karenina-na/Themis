package util

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"Themis/src/util/encryption"
	"Themis/src/util/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handle
// @Description: Handle the error and return the response
// @param        err : error
// @param        c   : gin.Context
func Handle(err error, c *gin.Context) {
	c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
	exception.HandleException(err)
}

// TokenError
//
//	@Description: Token error
//	@param c	: gin.Context
func TokenError(err error, c *gin.Context) {
	c.JSON(http.StatusUnauthorized, entity.NewFalseResult("false", "token-error"))
	exception.HandleException(err)
}

// CheckToken
//
//	@Description: Check the token
//	@param request	: *http.Request
//	@return Data	: *entity.Root
//	@return E	: error
func CheckToken(c *gin.Context) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckToken-util", util.Strval(r))
		}
	}()
	if !config.Root.TokenEnable {
		return true, nil
	}
	tokenStr := c.GetHeader("token")
	accountStr := c.GetHeader("account")
	passwordStr := c.GetHeader("password")
	account, password, err := token.ParseToken(tokenStr)
	if err != nil {
		return false, exception.NewUserError("CheckToken-util", "Token解析错误")
	}
	if account != encryption.Base64Decode(accountStr) || password != encryption.Base64Decode(passwordStr) {
		return false, exception.NewUserError("CheckToken-util", "token与用户数据不匹配")
	}
	if account != config.Root.RootAccount || password != config.Root.RootPassword {
		return false, exception.NewUserError("CheckToken-util", "账号验证失败")
	}
	return true, nil
}
