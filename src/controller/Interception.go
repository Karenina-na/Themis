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
		t := time.Now().UnixNano()
		defer func() {
			util.Loglevel(util.Debug, "Interception", "<-----------------------------------------------------")
			util.Loglevel(util.Info, "Interception", url)
			before := time.Now().UnixNano() - t
			util.Loglevel(util.Debug, "Interception", "耗时: "+
				strconv.FormatFloat(float64(before)/1000000000, 'f', 5, 64)+" s")
			util.Loglevel(util.Debug, "Interception", "----------------------------------------------------->")
		}()
		c.Next()
	}
}
