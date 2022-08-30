package Base

import (
	Init "Themis/src/Factory"
	"Themis/src/controller"
	"Themis/src/router"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FactoryBaseInit() *gin.Engine {
	arg := "debug"
	Init.ThemisInitFactory(&arg)
	r := gin.New()
	gin.SetMode(arg)
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFS("/static", http.Dir("./static"))
	r.Use(controller.Interception())
	router.MessageAPI(r)
	router.OperatorAPI(r)
	return r
}
