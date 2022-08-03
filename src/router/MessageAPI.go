package router

import (
	"Themis/src/controller"
	"github.com/gin-gonic/gin"
)

func MessageAPI(r *gin.Engine) {
	tx := r.Group("/api/message")
	{
		//	/api/message/register		服务注册
		tx.POST("/register", controller.RegisterController)
		//  /api/message/beat			心跳
		tx.PUT("/beat", controller.HeartBeat)
		//  /api/message/election		Leader调用开启选举
		tx.PUT("/election", controller.Election)
		// /api/message/getLeader		获取当前Leader
		tx.GET("/getLeader", controller.GetLeader)
		// /api/message/getServers		获取除Leader外其他服务
		tx.POST("/getServers", controller.GetServers)
	}
}
