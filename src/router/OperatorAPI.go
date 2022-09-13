package router

import (
	"Themis/src/controller"
	"Themis/src/controller/interception"
	"github.com/gin-gonic/gin"
)

// OperatorAPI
// @Description: OperatorAPI
// @param        r *gin.Engine
func OperatorAPI(r *gin.Engine) {
	tx := r.Group("/api/v1/operator")
	{
		//	/api/v1/operator/getInstances			获取全部服务实例
		tx.GET("/getInstances",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			controller.GetController)
		//  /api/v1/operator/getInstances			获取指定条件下的服务信息
		tx.POST("/getInstances",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			controller.GetPostController)
		// /api/v1/operator/deleteInstance			删除服务实例并拉入黑名单
		tx.DELETE("/deleteInstance",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			controller.DeleteInstanceController)
		// /api/v1/operator/deleteColony			删除地区集群实例并拉入黑名单
		tx.DELETE("/deleteColony",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			controller.DeleteColonyController)
		// /api/v1/operator/getDeleteInstance		获取被拉入黑名单的实例
		tx.GET("/getDeleteInstance",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			controller.GetDeleteInstanceController)
		// /api/v1/operator/cancelDeleteInstance	删除实例的黑名单
		tx.DELETE("/cancelDeleteInstance",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			controller.CancelDeleteInstanceController)
		// /api/v1/operator/getStatus				获取当前中心状态
		tx.GET("/getStatus", controller.GetStatusController)
		// /api/v1/operator/getColony				获取当前集群Leader
		tx.GET("/getClusterLeader", controller.GetClusterLeaderController)
		// /api/v1/operator/getColony				获取当前集群服务身份
		tx.GET("/getClusterStatus", controller.GetClusterStatusController)
	}
}
