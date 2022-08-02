package controller

import (
	"Envoy/src/entity"
	"Envoy/src/entity/util"
	"Envoy/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

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

func GetLeader(c *gin.Context) {
	c.JSON(http.StatusOK, entity.NewSuccessResult(service.GetLeader()))
}

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
