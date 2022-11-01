package router

import (
	"Themis/src/controller"
	"Themis/src/controller/interception"
	"github.com/gin-gonic/gin"
)

// MessageAPI
// @Description: MessageAPI
// @param        r *gin.Engine
func MessageAPI(r *gin.Engine) {
	//r.Use(interception.ClusterCandidateInterception()) // 拦截Candidate状态

	//允许Leader请求
	txLeader := r.Group("/v1/message/leader")
	txLeader.Use(interception.ClusterFollowInterception())
	{
		//	/v1/message/register		服务注册
		txLeader.POST("/register", controller.RegisterController)
		//  /v1/message/beat			心跳
		txLeader.PUT("/beat", controller.HeartBeatController)
		//  /v1/message/election		服务调用开启选举
		txLeader.PUT("/election", controller.ElectionController)
	}

	//允许Follow请求
	txFollow := r.Group("/v1/message/follow")
	txFollow.Use(interception.ClusterLeaderInterception())
	{
		// /v1/message/getLeader		获取当前Leader
		txFollow.POST("/getLeader", controller.GetLeaderController)
		// /v1/message/getServers		获取除Leader外其他服务
		txFollow.POST("/getServers", controller.GetServersController)
		// /v1/message/getServersNum	获取当前集群服务数量
		txFollow.POST("/getServersNum", controller.GetServersNumController)
	}
}
