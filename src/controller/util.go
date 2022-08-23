package controller

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Handle(err error, c *gin.Context) {
	c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
	exception.HandleException(err)
}
