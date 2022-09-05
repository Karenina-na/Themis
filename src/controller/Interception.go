package controller

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"github.com/gin-gonic/gin"
	"net/http"
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

// ClusterFollowInterception 拦截不属于follow的请求
func ClusterFollowInterception() gin.HandlerFunc {
	return func(c *gin.Context) {
		if syncBean.Status == syncBean.FOLLOW {
			c.Redirect(http.StatusTemporaryRedirect, "http://"+
				syncBean.LeaderAddress.IP+":"+syncBean.LeaderServicePort+
				c.Request.URL.Path)
			util.Loglevel(util.Info, "ClusterInterception", "重定向到Leader节点")
			c.Abort()
			return
		}
	}
}

// ClusterLeaderInterception 拦截不属于leader的请求
func ClusterLeaderInterception() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Cluster.ClusterEnable {
			if syncBean.Status == syncBean.LEADER {
				c.Next()
				return
			}

		}
	}
}

// ClusterCandidateInterception Candidate状态拒绝提供服务
func ClusterCandidateInterception() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Cluster.ClusterEnable {
			if syncBean.Status == syncBean.CANDIDATE {
				util.Loglevel(util.Info, "ClusterInterception", "当前节点为Candidate节点，无法处理请求")
				c.JSON(http.StatusOK, entity.NewFalseResult("false", "当前节点为Candidate节点，无法处理请求"))
				c.Abort()
				return
			}
		}
	}
}
