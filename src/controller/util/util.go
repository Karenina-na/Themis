package util

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
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

// CheckToken
//
//	@Description: Check the token
//	@param request	: *http.Request
//	@return Data	: *entity.Root
//	@return E	: error
func CheckToken(request *entity.RequestModel) (Data *interface{}, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckToken-util", util.Strval(r))
		}
	}()
	if !config.Root.TokenEnable {
		return &request.Data, nil
	}
	account, password, err := token.ParseToken(request.Root.Token)
	if err != nil {
		return nil, exception.NewUserError("CheckToken-util", "Token解析错误")
	}
	//if account != encryption.Base64Decode(request.Root.Account) ||
	//	password != encryption.Base64Decode(request.Root.Password) {
	//	return nil, exception.NewUserError("CheckToken-util", "token错误与用户数据不匹配")
	//}
	if account != request.Root.Account ||
		password != request.Root.Password {
		return nil, exception.NewUserError("CheckToken-util", "token错误与用户数据不匹配")
	}
	if account != config.Root.RootAccount || password != config.Root.RootPassword {
		return nil, exception.NewUserError("CheckToken-util", "token验证失败")
	}
	return &request.Data, nil
}
