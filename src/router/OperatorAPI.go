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

	// CURD
	txOperator := tx.Group("/CURD")
	txOperator.Use(interception.ClusterCandidateInterception(), interception.RootInterception())
	{
		//  /v1/operator/gerNamespaces	获取全部命名空间
		txOperator.GET("/getNamespaces", controller.GerNamespacesController)
		//  /v1/operator/getColonies	  获取指定命名空间下的集群
		txOperator.POST("/getColonies",
			interception.ClusterLeaderInterception(), controller.GetColoniesController)
		//  /v1/operator/getColoniesInstances	  获取指定命名空间下集群的节点
		txOperator.POST("/getColoniesInstances",
			interception.ClusterLeaderInterception(), controller.GetColoniesAndServerController)
		//	/v1/operator/getInstances			获取全部服务实例
		txOperator.GET("/getInstances",
			interception.ClusterLeaderInterception(), controller.GetAllServersController)
		//  /v1/operator/getInstancesByCondition			获取指定条件下的服务信息--模糊查询
		txOperator.POST("/getInstancesByCondition",
			interception.ClusterLeaderInterception(), controller.GetInstancesByConditionListController)
		// /v1/operator/getInstancesByConditionMap			获取指定条件下的服务信息--精确查询
		txOperator.POST("/getInstancesByConditionMap",
			interception.ClusterLeaderInterception(), controller.GetInstancesByConditionMapController)
		// /v1/operator/deleteInstance			删除服务实例并拉入黑名单
		txOperator.DELETE("/deleteInstance",
			interception.ClusterFollowInterception(), controller.DeleteInstanceController)
		// /v1/operator/deleteColony			删除地区集群实例并拉入黑名单
		txOperator.DELETE("/deleteColony",
			interception.ClusterFollowInterception(), controller.DeleteInstanceByColonyController)
		// /v1/operator/getDeleteInstance		获取被拉入黑名单的实例
		txOperator.GET("/getDeleteInstance",
			interception.ClusterLeaderInterception(), controller.GetDeleteInstancesController)
		// /v1/operator/cancelDeleteInstance	删除实例的黑名单
		txOperator.DELETE("/cancelDeleteInstance",
			interception.ClusterFollowInterception(), controller.CancelDeleteInstancesController)
	}

	//集群信息
	txCluster := r.Group("/v1/operator/cluster")
	txCluster.Use(interception.RootInterception())
	{
		// /v1/operator/getStatus				获取当前中心状态
		txCluster.GET("/getStatus", controller.GetStatusController)
		// /v1/operator/getClusterLeader				获取当前集群Leader
		txCluster.GET("/getClusterLeader", controller.GetClusterLeaderController)
		// /v1/operator/getClusterStatus				获取当前集群服务身份
		txCluster.GET("/getClusterStatus", controller.GetClusterStatusController)
	}

	//登录相关
	txLogin := r.Group("/v1/operator/manager")
	{
		// /v1/operator/manager/login				登录
		txLogin.POST("/login", controller.RootController)
	}
}
