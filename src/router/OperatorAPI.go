package router

import (
	"Themis/src/controller"
	"github.com/gin-gonic/gin"
)

func OperatorAPI(r *gin.Engine) {
	tx := r.Group("/api/v1/operator")
	{
		//	/api/v1/operator/getInstances			获取全部服务实例
		tx.GET("/getInstances", controller.GetController)
		// /api/v1/operator/deleteInstance			删除服务实例并拉入黑名单
		tx.DELETE("/deleteInstance", controller.DeleteInstanceController)
		// /api/v1/operator/deleteColony			删除地区集群实例并拉入黑名单
		tx.DELETE("/deleteColony", controller.DeleteColonyController)
		// /api/v1/operator/getDeleteInstances		获取被拉入黑名单的实例
		tx.GET("/getDeleteInstance", controller.GetDeleteInstanceController)
		// /api/v1/operator/cancelDeleteInstance	删除实例的黑名单
		tx.DELETE("/cancelDeleteInstance", controller.CancelDeleteInstanceController)
		// /api/v1/operator/getStatus				获取当前中心状态
		tx.GET("/getStatus", controller.GetStatusController)
	}
}
