package interception

import (
	"Themis/src/controller/util"
	"github.com/gin-gonic/gin"
)

// RootInterception
//
//	@Description: RootInterception	管理权限拦截器
//	@return gin.HandlerFunc	返回拦截器
func RootInterception() gin.HandlerFunc {
	return func(c *gin.Context) {
		b, err := util.CheckToken(c)
		if err != nil || !b {
			util.TokenError(err, c)
			c.Abort()
			return
		}
		c.Next()
	}
}
