package Base

import (
	Init "Themis/src/Factory"
	"Themis/src/controller/interception"
	"Themis/src/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

// FactoryBaseInit
// @Description: 初始化基础模块
// @return       *gin.Engine 返回路由引擎
func FactoryBaseInit() *gin.Engine {
	arg := "debug"
	Init.ThemisInitFactory(&arg)
	r := gin.New()
	gin.SetMode(arg)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFS("/static", http.Dir("./static"))
	r.Use(interception.Interception())
	router.MessageAPI(r)
	router.OperatorAPI(r)
	return r
}
