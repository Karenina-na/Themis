package main

import (
	Init "Themis/src/Factory"
	"Themis/src/config"
	"Themis/src/controller/interception"
	"Themis/src/util"
	"flag"
	swaggerFile "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

import (
	"Themis/src/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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
*/

// @title          Themis
// @version        1.0
// @description    分布式记账系统调度中心
// @termsOfService https://www.wzxaugenstern.online/#/
// @contact.name   CYCLEWW
// @contact.url    https://www.wzxaugenstern.online/#/
// @contact.email  1539989223@qq.com
// @license.name   Apache 2.0
// @license.url    http://www.apache.org/licenses/LICENSE-2.0.html
//
// main
// @Description:   主函数
func main() {
	arg := flag.String("mode", "debug", "debug / release /test 环境")
	flag.Parse()
	Init.ThemisInitFactory(arg)
	if *arg == "debug" {
		gin.SetMode(gin.DebugMode)
		util.Loglevel(util.Info, "main", "debug mode")
	} else if *arg == "release" {
		gin.SetMode(gin.ReleaseMode)
		util.Loglevel(util.Info, "main", "release mode")
	} else if *arg == "test" {
		gin.SetMode(gin.TestMode)
		util.Loglevel(util.Info, "main", "test mode")
	}
	defer func() {
		err := recover()
		if err != nil {
			util.Loglevel(util.Error, "main", util.Strval(err))
		}
	}()
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(cors.Default())
	r.Use(interception.Interception())
	if *arg == "debug" {
		r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFile.Handler))
	}
	router.MessageAPI(r)
	router.OperatorAPI(r)
	go func() {
		util.Loglevel(util.Info, "main", "start service on port:"+config.Port.CenterPort)
		err := r.Run(":" + config.Port.CenterPort)
		if err != nil {
			util.Loglevel(util.Error, "main", util.Strval(err))
			util.Loglevel(util.Error, "main", "server start error")
			os.Exit(0)
		}
	}()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	util.Loglevel(util.Info, "main", "Themis is exiting...")
	Init.ThemisCloseFactory()
	runtime.GC()
	util.Loglevel(util.Info, "main", "Themis is exited")
	time.Sleep(time.Second * 3)
}
