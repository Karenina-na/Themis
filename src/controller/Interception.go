package controller

import (
	"Themis/src/util"
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

func Interception() gin.HandlerFunc {
	return func(c *gin.Context) {
		url := c.Request.Host + c.Request.URL.Path
		util.Loglevel(util.Debug, "Interception", "<-----------------------------------------------------")
		t := time.Now().Unix()
		util.Loglevel(util.Info, "Interception", url)
		c.Next()
		util.Loglevel(util.Debug, "Interception", "耗时-"+strconv.FormatInt(time.Now().Unix()-t, 10))
		util.Loglevel(util.Debug, "Interception", "----------------------------------------------------->")
	}
}
