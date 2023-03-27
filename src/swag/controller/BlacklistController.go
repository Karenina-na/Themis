package controller

import (
	"Themis/src/controller/util"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// BlackInstanceController
// @Summary     删除服务实例并拉入黑名单
// @Description 由管理者调用删除服务实例并拉入黑名单。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "被删除的服务实例信息"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/CURD/blacklistInstance [delete]
func BlackInstanceController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("BlackInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		util.Handle(err1, c)
		return
	}
	if Assert1 {
		Assert2, err2 := service.BlackInstance(Server)
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

// BlackInstanceByColonyController
// @Summary     删除地区集群实例并拉入黑名单
// @Description 由管理者调用删除地区集群实例并拉入黑名单。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "被删除的服务地区信息（用服务实例信息包装）"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/CURD/blacklistColony [delete]
func BlackInstanceByColonyController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("BlackInstanceByColonyController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert, err1 := service.BlackColony(Server)
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

// BlackInstanceByNamespaceController
// @Summary     删除命名空间集群实例并拉入黑名单
// @Description 由管理者调用删除地区集群实例并拉入黑名单。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "被删除的服务命名空间信息（用服务实例信息包装）"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/CURD/blacklistNamespace [delete]
func BlackInstanceByNamespaceController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("BlackInstanceByColonyController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert, err1 := service.BlackNamespace(Server)
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

// GetBlacklistInstancesController
// @Summary     获取全部黑名单服务实例
// @Description 由管理者调用的获取当前全部黑名单服务实例
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{data=[]entity.ServerModel} "返回黑名单中服务实例切片数组"
// @Router      /operator/CURD/getBlacklist [GET]
func GetBlacklistInstancesController(c *gin.Context) {
	servers, err := service.GetBlacklistServer()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// DeleteBlacklistInstancesController
// @Summary     删除黑名单中的实例信息
// @Description 由管理者调用删除黑名单中的实例信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel true "从黑名单中清除的实例信息"
// @Success     200   {object} entity.ResultModel "返回true或false"
// @Router      /operator/CURD/deleteBlacklistInstance [delete]
func DeleteBlacklistInstancesController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("DeleteBlacklistInstancesController", "参数绑定错误-"+err.Error()))
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
