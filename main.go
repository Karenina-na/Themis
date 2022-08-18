package main

import (
	"Themis/docs"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

import (
	"Themis/src/config"
	"Themis/src/controller"
	"Themis/src/entity/util"
	"Themis/src/router"
	"github.com/gin-gonic/gin"
) // gin-swagger middleware

// @title Themis API
// @version 1.0
// @description 分布式记账系统调度中心
// @termsOfService https://www.wzxaugenstern.online/#/
// @contact.name CYCLEWW
// @contact.url https://www.wzxaugenstern.online/#/
// @contact.email 1539989223@qq.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	config.SwaggerConfig(docs.SwaggerInfo)
	defer func() {
		err := recover()
		util.Loglevel(util.Error, "main", util.Strval(err))
	}()
	r := gin.Default()
	r.Use(controller.Interception())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	router.MessageAPI(r)
	router.OperatorAPI(r)
	err := r.Run(":" + config.Port)
	if err != nil {
		util.Loglevel(util.Error, "main", util.Strval(err))
		panic(err.Error())
		return
	}
}
