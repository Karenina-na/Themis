package router

import (
	"Themis/src/controller"
	"Themis/src/controller/interception"
	"github.com/gin-gonic/gin"
)

//
// MessageAPI
// @Description: MessageAPI
// @param        r *gin.Engine
//
func MessageAPI(r *gin.Engine) {
	tx := r.Group("/api/v1/message")
	tx.Use(interception.ClusterCandidateInterception())
	{
		//	/api/v1/message/register		服务注册
		tx.POST("/register", interception.ClusterFollowInterception(),
			controller.RegisterController)
		//  /api/v1/message/beat			心跳
		tx.PUT("/beat", interception.ClusterFollowInterception(),
			controller.HeartBeatController)
		//  /api/v1/message/election		服务调用开启选举
		tx.PUT("/election", interception.ClusterFollowInterception(),
			controller.ElectionController)
		// /api/v1/message/getLeader		获取当前Leader
		tx.POST("/getLeader", interception.ClusterLeaderInterception(),
			controller.GetLeaderController)
		// /api/v1/message/getServers		获取除Leader外其他服务
		tx.POST("/getServers", interception.ClusterLeaderInterception(),
			controller.GetServersController)
		// /api/v1/message/getServersNum	获取当前集群服务数量
		tx.POST("/getServersNum", interception.ClusterLeaderInterception(),
			controller.GetServersNumController)
	}
}
