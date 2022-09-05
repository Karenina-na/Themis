package router

import (
	"Themis/src/controller"
	"github.com/gin-gonic/gin"
)

func OperatorAPI(r *gin.Engine) {
	tx := r.Group("/api/v1/operator")
	{
		//	/api/v1/operator/getInstances			获取全部服务实例
		tx.GET("/getInstances", controller.ClusterLeaderInterception(),
			controller.GetController)
		//  /api/v1/operator/getInstances			获取指定条件下的服务信息
		tx.POST("/getInstances", controller.ClusterLeaderInterception(),
			controller.GetPostController)
		// /api/v1/operator/deleteInstance			删除服务实例并拉入黑名单
		tx.DELETE("/deleteInstance", controller.ClusterFollowInterception(),
			controller.DeleteInstanceController)
		// /api/v1/operator/deleteColony			删除地区集群实例并拉入黑名单
		tx.DELETE("/deleteColony", controller.ClusterFollowInterception(),
			controller.DeleteColonyController)
		// /api/v1/operator/getDeleteInstance		获取被拉入黑名单的实例
		tx.GET("/getDeleteInstance", controller.ClusterLeaderInterception(),
			controller.GetDeleteInstanceController)
		// /api/v1/operator/cancelDeleteInstance	删除实例的黑名单
		tx.DELETE("/cancelDeleteInstance", controller.ClusterFollowInterception(),
			controller.CancelDeleteInstanceController)
		// /api/v1/operator/getStatus				获取当前中心状态
		tx.GET("/getStatus", controller.GetStatusController)
	}
}
