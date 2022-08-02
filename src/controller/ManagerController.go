package controller

import (
	"Envoy/src/entity"
	"Envoy/src/entity/util"
	"Envoy/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetController(c *gin.Context) {
	servers := service.GetInstances()
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}
func GetDeleteInstanceController(c *gin.Context) {
	servers := service.GetDeleteInstances()
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

func DeleteInstanceController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		util.Loglevel(util.Warn, "DeleteInstanceController-Controller", "参数绑定错误-"+err.Error())
	} else {
		if service.CheckServer(Server) {
			Assert := service.DeleteServer(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("false", "实例不存在"))
		}
	}
}

func DeleteColonyController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		util.Loglevel(util.Warn, "DeleteColonyController-Controller", "参数绑定错误-"+err.Error())
	} else {
		Assert := service.DeleteColony(Server)
		c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
	}
}

func CancelDeleteInstanceController(c *gin.Context) {
	Server := entity.NewServerModel()
	err := c.BindJSON(Server)
	if err != nil {
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		util.Loglevel(util.Warn, "CancelDeleteInstanceController-Controller", "参数绑定错误-"+err.Error())
	} else {
		if service.CheckDeleteServer(Server) {
			Assert := service.DeleteDeleteInstance(Server)
			c.JSON(http.StatusOK, entity.NewSuccessResult(Assert))
		} else {
			c.JSON(http.StatusOK, entity.NewFalseResult("false", "实例不在拒绝队列中"))
		}
	}
}
