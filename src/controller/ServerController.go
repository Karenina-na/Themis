package controller

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// RegisterController
// @Summary 服务注册
// @Description 服务注册进中心
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/message/register [post]
func RegisterController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("RegisterController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	assert2, err2 := service.CheckDeleteServer(Server)
	if err2 != nil {
		handle(err2)
		return
	}
	if !assert1 && !assert2 {
		Assert, err3 := service.RegisterServer(Server)
		if err3 != nil {
			handle(err3)
			return
		} else {
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		}
	} else if assert1 && !assert2 {
		c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例已注册"))
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例已被删除"))
	}

}

// HeartBeatController
// @Summary 服务心跳
// @Description 服务心跳重置倒计时
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/message/message/beat [put]
func HeartBeatController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("HeartBeatController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	Assert2, err2 := service.CheckDeleteServer(Server)
	if err2 != nil {
		handle(err2)
		return
	}
	if Assert1 && !Assert2 {
		Assert3, err3 := service.FlashHeartBeat(Server)
		if err3 != nil {
			handle(err3)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert3))
	} else if !Assert2 {
		Assert4, err4 := service.RegisterServer(Server)
		if err4 != nil {
			handle(err4)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert4))
	} else {
		c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例已被删除"))
	}
}

// ElectionController
// @Summary 选举
// @Description 由领导者调用的新一轮选举接口。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "领导者服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/message/election [put]
func ElectionController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("ElectionController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckServer(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	if Assert1 {
		Assert2, err2 := service.Election(Server)
		if err2 != nil {
			handle(err2)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert2))
	} else {
		exception.HandleException(exception.NewControllerPanic("ElectionController", "错误的Leader"))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "错误的Leader"))
	}
}

// GetLeaderController
// @Summary 获取领导者
// @Description 由其他服务调用的获取当前领导者接口。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回领导者服务信息"
// @Router /api/v1/message/getLeader [GET]
func GetLeaderController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	leader, err := service.GetLeader()
	if err != nil {
		handle(err)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(leader))
}

// GetServersController
// @Summary 获取当前被领导者服务列表
// @Description 由当前领导者调用的获取领导者所领导的服务列表。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回被领导者的切片数组"
// @Router /api/v1/message/getServers [POST]
func GetServersController(c *gin.Context) {
	handle := func(err any) {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "服务端异常"))
		exception.HandleException(err)
	}
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		exception.HandleException(exception.NewControllerPanic("GetServersController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	Assert1, err1 := service.CheckLeader(Server)
	if err1 != nil {
		handle(err1)
		return
	}
	if Assert1 {
		Assert2, err2 := service.GetServers(Server)
		if err2 != nil {
			handle(err2)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert2))
	} else {
		exception.HandleException(exception.NewControllerPanic("GetServersController", "错误的Leader"))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "错误的Leader"))
	}
}
