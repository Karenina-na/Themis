package controller

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetController
// @Summary 获取全部服务实例
// @Description 由管理者调用的获取当前所有服务信息。
// @Tags 管理层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回服务实例切片数组"
// @Router /api/v1/operator/getInstances [get]
func GetController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	servers, err := service.GetInstances()
	if err != nil {
		handle(err)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// DeleteInstanceController
// @Summary 删除服务实例并拉入黑名单
// @Description 由管理者调用删除服务实例并拉入黑名单。
// @Tags 管理层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "被删除的服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/operator/election [delete]
func DeleteInstanceController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("DeleteInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	if Assert1 {
		Assert2, err2 := service.DeleteServer(Server)
		if err2 != nil {
			handle(err2)
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
// @Summary 删除地区集群实例并拉入黑名单
// @Description 由管理者调用删除地区集群实例并拉入黑名单。
// @Tags 管理层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "被删除的服务地区信息（用服务实例信息包装）"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/operator/deleteColony [delete]
func DeleteColonyController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("DeleteColonyController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert, err1 := service.DeleteColony(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))

}

// GetDeleteInstanceController
// @Summary 获取全部黑名单服务实例
// @Description 由管理者调用的获取当前全部黑名单服务实例
// @Tags 管理层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回黑名单中服务实例切片数组"
// @Router /api/v1/operator/getDeleteInstance [get]
func GetDeleteInstanceController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	servers, err := service.GetDeleteInstances()
	if err != nil {
		handle(err)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// CancelDeleteInstanceController
// @Summary 删除黑名单中的实例信息
// @Description 由管理者调用删除黑名单中的实例信息。
// @Tags 管理层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "从黑名单中清除的实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/operator/cancelDeleteInstance [delete]
func CancelDeleteInstanceController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("CancelDeleteInstanceController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckDeleteServer(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	if Assert1 {
		Assert2, err2 := service.DeleteDeleteInstance(Server)
		if err2 != nil {
			handle(err2)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert2))
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "实例不在拒绝队列中"))
	}
}
