package controller

import (
	"Themis/src/controller/util"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GerNamespacesController
// @Summary     获取全部命名空间
// @Description 由管理者调用的获取当前所有服务信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{data=[]string} "返回命名空间列表名称"
// @Router      /operator/CURD/getNamespaces [GET]
func GerNamespacesController(c *gin.Context) {
	namespaces, err := service.GetNamespaces()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(namespaces))
}

// GetColoniesController
// @Summary     获取指定命名空间下的集群
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{data=[]string} "返回集群列表名称"
// @Router      /operator/CURD/getColonies [POST]
func GetColoniesController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("GetColoniesController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	colony, err := service.GetColonyByNamespace(Server.Namespace)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(colony))
}

// GetColoniesAndServerController
// @Summary     获取指定命名空间下的集群和服务名
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{data=map[string][]string} "返回集群服务列表名称"
// @Router      /operator/CURD/getColoniesInstances [POST]
func GetColoniesAndServerController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("GetColoniesAndServerController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	colony, err := service.GetColonyAndServerByNamespace(Server.Namespace)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(colony))
}

// GetAllServersController
// @Summary     获取全部服务实例
// @Description 由管理者调用的获取当前所有服务信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{} "返回服务全体封装"
// @Router      /operator/CURD/getInstances [GET]
func GetAllServersController(c *gin.Context) {
	servers, err := service.GetInstances()
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// GetServerController
// @Summary     获取服务实例
// @Description 由管理者调用的获取指定服务信息。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Success     200   {object} entity.ResultModel{} "返回服务封装"
// @Router      /operator/CURD/getInstance [POST]
func GetServerController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("GetInstancesByConditionListController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	servers, err := service.GetInstance(Server)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(servers))
}

// GetInstancesByConditionListController
// @Summary     获取指定条件下的服务实例服务实例--模糊查询
// @Description 获取命名空间与集群条件下的服务实例。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel                            true "封装的条件参数"
// @Success     200 {object} entity.ResultModel{data=[]entity.ServerModel} "返回服务实例切片数组"
// @Router      /operator/CURD/getInstancesByCondition [POST]
func GetInstancesByConditionListController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("GetInstancesByConditionListController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	list, err := service.GetInstancesByNamespaceAndColony(Server)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(list))
}

// GetInstancesByConditionMapController
// @Summary     获取指定条件下的服务实例服务实例--精确查询
// @Description 获取命名空间与集群条件下的服务实例。
// @Tags        管理层
// @Accept      application/json
// @Produce     application/json
// @Security    ApiKeyAuth
// @Param       Model body     entity.ServerModel                            true "封装的条件参数"
// @Success     200 {object} entity.ResultModel{data=map[string][]entity.ServerModel} "返回服务实例map切片数组"
// @Router      /operator/CURD/getInstancesByConditionMap [POST]
func GetInstancesByConditionMapController(c *gin.Context) {
	Server := entity.NewServerModel()
	if err := c.BindJSON(Server); err != nil {
		exception.HandleException(exception.NewUserError("GetInstancesByConditionMapController", "参数绑定错误-"+err.Error()))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-"+err.Error()))
		return
	}
	if Server.Namespace == "" || Server.Colony == "" {
		exception.HandleException(exception.NewUserError("GetInstancesByConditionMapController", "参数绑定错误-命名空间或集群为空"))
		c.JSON(http.StatusOK, entity.NewFalseResult("false", "参数绑定错误-命名空间或集群为空"))
		return
	}
	list, err := service.GetInstanceByNamespaceAndColony(Server)
	if err != nil {
		util.Handle(err, c)
		return
	}
	c.JSON(http.StatusOK, entity.NewSuccessResult(list))
}
