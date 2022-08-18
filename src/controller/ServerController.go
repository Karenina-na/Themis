package controller

import (
	"Themis/src/entity"
	"Themis/src/entity/util"
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
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		util.Loglevel(util.Warn, "RegisterController-Controller", "参数绑定错误-"+err.Error())
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
	} else {
		if !service.CheckServer(Server) && !service.CheckDeleteServer(Server) {
			Assert := service.RegisterServer(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else if service.CheckServer(Server) && !service.CheckDeleteServer(Server) {
			c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例已注册"))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例已被删除"))
		}
	}
}

// HeartBeat
// @Summary 服务心跳
// @Description 服务心跳重置倒计时
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/message/message/beat [put]
func HeartBeat(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		util.Loglevel(util.Warn, "HeartBeat-Controller", "参数绑定错误-"+err.Error())
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
	} else {
		if service.CheckServer(Server) && !service.CheckDeleteServer(Server) {
			Assert := service.FlashHeartBeat(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else if !service.CheckDeleteServer(Server) {
			Assert := service.RegisterServer(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("False", "实例以被删除"))
		}
	}
}

// Election
// @Summary 选举
// @Description 由领导者调用的新一轮选举接口。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Param Model query entity.ServerModel true "领导者服务实例信息"
// @Success 200 {object} entity.ResultModel "返回true或false"
// @Router /api/v1/message/election [put]
func Election(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		util.Loglevel(util.Warn, "Election-Controller", "参数绑定错误-"+err.Error())
	} else {
		if service.CheckLeader(Server) {
			Assert := service.Election(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("false", "错误的Leader"))
			util.Loglevel(util.Warn, "Election-Controller", "错误的Leader")
		}
	}
}

// GetLeader
// @Summary 获取领导者
// @Description 由其他服务调用的获取当前领导者接口。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回领导者服务信息"
// @Router /api/v1/message/getLeader [GET]
func GetLeader(c *gin.Context) {
	c.JSON(http.StatusOK, entity.NewSuccessResult(service.GetLeader()))
}

// GetServers
// @Summary 获取当前被领导者服务列表
// @Description 由当前领导者调用的获取领导者所领导的服务列表。
// @Tags 服务层
// @Accept application/json
// @Produce application/json
// @Security ApiKeyAuth
// @Success 200 {object} entity.ResultModel "返回被领导者的切片数组"
// @Router /api/v1/message/getServers [POST]
func GetServers(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(&Server)
	if err != nil {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		util.Loglevel(util.Warn, "GetLeader-Controller", "参数绑定错误-"+err.Error())
	} else {
		if service.CheckLeader(Server) {
			Assert := service.GetServers(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("false", "错误的Leader"))
			util.Loglevel(util.Error, "GetLeader-Controller", "错误的Leader")
		}
	}
}
