package router

import (
	"Themis/src/controller"
	"github.com/gin-gonic/gin"
)

func MessageAPI(r *gin.Engine) {
	tx := r.Group("/api/v1/message")
	{
		//	/api/v1/message/register		服务注册
		tx.POST("/register", controller.ClusterFollowInterception(),
			controller.RegisterController)
		//  /api/v1/message/beat			心跳
		tx.PUT("/beat", controller.ClusterFollowInterception(),
			controller.HeartBeatController)
		//  /api/v1/message/election		服务调用开启选举
		tx.PUT("/election", controller.ClusterFollowInterception(),
			controller.ElectionController)
		// /api/v1/message/getLeader		获取当前Leader
		tx.POST("/getLeader", controller.ClusterLeaderInterception(),
			controller.GetLeaderController)
		// /api/v1/message/getServers		获取除Leader外其他服务
		tx.POST("/getServers", controller.ClusterLeaderInterception(),
			controller.GetServersController)
		// /api/v1/message/getServersNum	获取当前集群服务数量
		tx.POST("/getServersNum", controller.ClusterLeaderInterception(),
			controller.GetServersNumController)
	}
}
