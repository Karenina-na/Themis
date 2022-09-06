package util

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"github.com/gin-gonic/gin"
	"net/http"
)

//
// Handle
// @Description: Handle the error and return the response
// @param        err : error
// @param        c   : gin.Context
//
func Handle(err error, c *gin.Context) {
	c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
	exception.HandleException(err)
}
