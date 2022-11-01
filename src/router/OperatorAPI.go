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
	tx := r.Group("/v1/operator")
	{
		//  /v1/operator/gerNamespaces	获取全部命名空间
		tx.GET("/gerNamespaces",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.GerNamespacesController)
		//  /v1/operator/getColonys	  获取指定命名空间下的集群
		tx.GET("/getColonys",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.GetColonyController)
		//	/v1/operator/getInstances			获取全部服务实例
		tx.GET("/getInstances",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.GetController)
		//  /v1/operator/getInstancesByCondition			获取指定条件下的服务信息
		tx.POST("/getInstancesByCondition",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.GetPostController)
		// /v1/operator/deleteInstance			删除服务实例并拉入黑名单
		tx.DELETE("/deleteInstance",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.DeleteInstanceController)
		// /v1/operator/deleteColony			删除地区集群实例并拉入黑名单
		tx.DELETE("/deleteColony",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.DeleteColonyController)
		// /v1/operator/getDeleteInstance		获取被拉入黑名单的实例
		tx.GET("/getDeleteInstance",
			interception.ClusterLeaderInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.GetDeleteInstanceController)
		// /v1/operator/cancelDeleteInstance	删除实例的黑名单
		tx.DELETE("/cancelDeleteInstance",
			interception.ClusterFollowInterception(), interception.ClusterCandidateInterception(),
			interception.RootInterception(),
			controller.CancelDeleteInstanceController)
		// /v1/operator/getStatus				获取当前中心状态
		tx.GET("/getStatus", interception.RootInterception(),
			controller.GetStatusController)
		// /v1/operator/getClusterLeader				获取当前集群Leader
		tx.GET("/getClusterLeader", interception.RootInterception(),
			controller.GetClusterLeaderController)
		// /v1/operator/getClusterStatus				获取当前集群服务身份
		tx.GET("/getClusterStatus", interception.RootInterception(),
			controller.GetClusterStatusController)
		// /v1/operator/login				登录
		tx.POST("/login", controller.RootController)
	}
}
