package controller

import (
	"Themis/src/config"
	"Themis/src/controller/util"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"Themis/src/util/encryption"
	"Themis/src/util/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetStatusController
// @Summary     获取服务器状态
// @Description 由管理员调用获取当前中心服务器状态
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200 {object} entity.ResultModel{data=entity.ComputerInfoModel} "返回电脑状态"
// @Router      /operator/cluster/getStatus [GET]
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
// @Router      /operator/cluster/getClusterLeader [GET]
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
// @Router      /operator/cluster/getClusterStatus [GET]
func GetClusterStatusController(c *gin.Context) {
	leader, err := service.GetClusterStatus()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(leader))
}

// RootController
// @Summary     管理员登录
// @Description 管理员登录接口。
// @Tags        服务层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.Root true "管理员账号密码"
// @Success     200   {object} entity.ResultModel "返回token"
// @Router      /message/manager/login [POST]
func RootController(c *gin.Context) {
	root := entity.NewRootModel()
	err := c.BindJSON(&root)
	if err != nil {
		exception.HandleException(exception.NewUserError("GetServersController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	account := encryption.Base64Decode(root.Account)
	password := encryption.Base64Decode(root.Password)
	if (account == config.Root.RootAccount) && (password == config.Root.RootPassword) {
		t, err := token.GenerateToken(config.Root.RootAccount, config.Root.RootPassword)
		if err != nil {
			util.Handle(err, c)
			return
		}
		c.JSON(http.StatusOK, entity.NewSuccessResult(struct {
			Token string `json:"token"`
		}{Token: t}))
		return
	}
	c.JSON(http.StatusOK, entity.NewFalseResult("false", "密码错误"))
}
