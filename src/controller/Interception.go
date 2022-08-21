package controller

import (
	"github.com/gin-gonic/gin"
)

func Interception() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
