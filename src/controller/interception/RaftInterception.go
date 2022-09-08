package interception

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

//
// ClusterFollowInterception
// @Description: 拦截不属于Follow的请求
// @return       gin.HandlerFunc 返回拦截器
//
func ClusterFollowInterception() gin.HandlerFunc {
	return func(c *gin.Context) {
		if syncBean.Status == syncBean.FOLLOW {
			c.Redirect(http.StatusTemporaryRedirect, "http://"+
				syncBean.Leader.LeaderAddress.IP+":"+syncBean.Leader.LeaderAddress.Port+
				c.Request.URL.Path)
			util.Loglevel(util.Info, "ClusterInterception", "重定向到Leader节点")
			c.Abort()
			return
		}
	}
}

//
// ClusterLeaderInterception
// @Description: 拦截不属于Leader的请求
// @return       gin.HandlerFunc 返回拦截器
//
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

//
// ClusterCandidateInterception
// @Description: Candidate节点拦截器
// @return       gin.HandlerFunc 返回拦截器
//
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
