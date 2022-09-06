package controller

import (
	"Themis/src/controller/util"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetController
// @Summary     获取全部服务实例
// @Description 由管理者调用的获取当前所有服务信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{data=[]entity.ServerModel} "返回服务实例切片数组"
// @Router      /operator/getInstances [get]
func GetController(c *gin.Context) {
	servers, err := service.GetInstances()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// GetPostController
// @Summary     获取指定条件下的服务实例服务实例
// @Description 由获取命名空间与集群条件下的服务实例服务实例。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel                            true "封装的条件参数"
// @Success     200 {object} entity.ResultModel{data=[]entity.ServerModel} "返回服务实例切片数组"
// @Router      /operator/getInstances [POST]
func GetPostController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewUserError("DeleteInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	namespaces, err := service.GetInstancesByNamespaceAndColony(Server)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(namespaces))
}

// DeleteInstanceController
// @Summary     删除服务实例并拉入黑名单
// @Description 由管理者调用删除服务实例并拉入黑名单。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "被删除的服务实例信息"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/deleteInstance [delete]
func DeleteInstanceController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewUserError("DeleteInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		util.Handle(err1, c)
		return
	}
	if Assert1 {
		Assert2, err2 := service.DeleteServer(Server)
		if err2 != nil {
			util.Handle(err2, c)
			return
		}
		if Assert2 {
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert2))
		}
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "实例不存在"))
	}
}

// DeleteColonyController
// @Summary     删除地区集群实例并拉入黑名单
// @Description 由管理者调用删除地区集群实例并拉入黑名单。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "被删除的服务地区信息（用服务实例信息包装）"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/deleteColony [delete]
func DeleteColonyController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewUserError("DeleteColonyController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert, err1 := service.DeleteColonyServer(Server)
	if err1 != nil {
		util.Handle(err1, c)
		return
	}
	if Assert {
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "命名空间或集群不存在"))
	}
}

// GetDeleteInstanceController
// @Summary     获取全部黑名单服务实例
// @Description 由管理者调用的获取当前全部黑名单服务实例
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{data=[]entity.ServerModel} "返回黑名单中服务实例切片数组"
// @Router      /operator/getDeleteInstance [get]
func GetDeleteInstanceController(c *gin.Context) {
	servers, err := service.GetBlacklistServer()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// CancelDeleteInstanceController
// @Summary     删除黑名单中的实例信息
// @Description 由管理者调用删除黑名单中的实例信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "从黑名单中清除的实例信息"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/cancelDeleteInstance [delete]
func CancelDeleteInstanceController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewUserError("CancelDeleteInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckDeleteServer(Server)
	if err1 != nil {
		util.Handle(err1, c)
		return
	}
	if Assert1 {
		Assert2, err2 := service.DeleteInstanceFromBlacklist(Server)
		if err2 != nil {
			util.Handle(err2, c)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert2))
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "实例不在拒绝队列中"))
	}
}

// GetStatusController
// @Summary     获取服务状态
// @Description 由管理员调用获取当前中心线程数
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{data=entity.ComputerInfoModel} "返回电脑状态"
// @Router      /operator/getStatus [get]
func GetStatusController(c *gin.Context) {
	computer, err := service.GetCenterStatus()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(computer))
}

// GetClusterLeaderController
// @Summary     获取中心集群领导者
// @Description 由管理员调用获取中心集群领导者
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{} "返回中心集群领导者名称"
// @Router      /operator/getClusterLeader [get]
func GetClusterLeaderController(c *gin.Context) {
	leader, err := service.GetClusterLeader()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(leader))
}

// GetClusterStatusController
// @Summary     获取中心当前身份状态
// @Description 由管理员调用获取中心当前身份状态
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{data=syncBean.StatusLevel} "返回集群状态"
// @Router      /operator/getClusterStatus [get]
func GetClusterStatusController(c *gin.Context) {
	leader, err := service.GetClusterStatus()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(leader))
}
