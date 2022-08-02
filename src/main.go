package main

import (
	"Envoy/src/config"
	"Envoy/src/controller"
	"Envoy/src/entity/util"
	"Envoy/src/router"
	"github.com/gin-gonic/gin"
)

func main() {
	defer func() {
		err := recover()
		util.Loglevel(util.Error, "main", util.Strval(err))
	}()
	r := gin.Default()
	r.Use(controller.Interception())
	router.MessageAPI(r)
	router.OperatorAPI(r)
	err := r.Run(":" + config.Port)
	if err != nil {
		util.Loglevel(util.Error, "main", util.Strval(err))
		panic(err.Error())
		return
	}
}
