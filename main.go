package main

import (
	FactoryInit "Themis/src/Init"
	"Themis/src/config"
	util2 "Themis/src/util"
	"flag"
	"fmt"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

import (
	"Themis/src/controller"
	"Themis/src/router"
	"github.com/gin-gonic/gin"
) // gin-swagger middleware

/*
~ Licensed to the Apache Software Foundation (ASF) under one or more
~ contributor license agreements.  See the NOTICE file distributed with
~ this work for additional information regarding copyright ownership.
~ The ASF licenses this file to You under the Apache License, Version 2.0
~ (the "License"); you may not use this file except in compliance with
~ the License.  You may obtain a copy of the License at
~
~     http://www.apache.org/licenses/LICENSE-2.0
~
~ Unless required by applicable law or agreed to in writing, software
~ distributed under the License is distributed on an "AS IS" BASIS,
~ WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
~ See the License for the specific language governing permissions and
~ limitations under the License.
~
*/

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
	arg := flag.String("mode", "debug", "debug / release /test 环境")
	flag.Parse()
	if *arg == "debug" {
		gin.SetMode(gin.DebugMode)
		fmt.Println("debug mode")
	} else if *arg == "release" {
		gin.SetMode(gin.ReleaseMode)
		fmt.Println("release mode")
	} else if *arg == "test" {
		gin.SetMode(gin.TestMode)
		fmt.Println("test mode")
	}
	FactoryInit.ThemisInitFactory()
	defer func() {
		err := recover()
		util2.Loglevel(util2.Error, "main", util2.Strval(err))
	}()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.StaticFS("/static", http.Dir("./static"))
	r.Use(controller.Interception())
	if *arg == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	}
	router.MessageAPI(r)
	router.OperatorAPI(r)
	go func() {
		err := r.Run(":" + config.Port)
		if err != nil {
			util2.Loglevel(util2.Error, "main", util2.Strval(err))
			os.Exit(0)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util2.Loglevel(util2.Info, "main", "Themis is exiting...")
	runtime.GC()
	time.Sleep(1 * time.Second)
}
