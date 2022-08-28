package main

import (
	"Themis/test/Base"
	"Themis/test/ServerTest"
	"fmt"
	"github.com/gin-gonic/gin"
	"testing"
	"time"
)

/*
model1:= ip:111.111.111.111 namespace:default
model2:= ip:222.222.222.222 namespace:default
model3:= ip:333.333.333.333 namespace:A
model4:= ip:444.444.444.444 namespace:A
*/
func Test(t *testing.T) {
	//ServerTest.RegisterTest(router)	// 注册
	//ServerTest.ElectionTest(router)	// 选举
	//ServerTest.GetLeaderTest(router)	// 获取leader实例
	//ServerTest.CancelDeleteInstance(router)	// 取消删除实例
	//ServerTest.DeleteServerTest(router)	// 删除实例
	//ServerTest.GetAllServersTest(router)	// 获取所有实例
	//ServerTest.GetDeleteInstances(router)	// 获取删除实例
	//ServerTest.GetFollowServer(router)	// 获取跟随实例
	//ServerTest.GetInstancesByQuery(router)	// 根据条件获取实例
	//ServerTest.DeleteServerByColony(router)	// 根据条件删除实例
	//ServerTest.ServerBeatTest(router)	// 服务心跳
	router := Base.FactoryBaseInit()
	time.Sleep(time.Second * 2)

	/*获取电脑状态*/
	ServerTest.GetComputerStatus(router) // 获取电脑状态

	/*注册-选举-拉入黑名单-拉出黑名单*/
	test1(router)

	/*注册-选举-获取leader领导的实例-拉入黑名单-拉出黑名单*/
	test2(router)

	/*注册-获取指定条件下服务器列表-拉入黑名单-拉出黑名单*/
	test3(router)

	/*注册-心跳-拉入黑名单-拉出黑名单*/
	test4(router)
}

/**
注册-选举-拉入黑名单-拉出黑名单
*/
func test1(router *gin.Engine) {
	ServerTest.RegisterTest(router)
	ServerTest.GetAllServersTest(router)
	ServerTest.ElectionTest(router)
	ServerTest.GetLeaderTest(router)
	ServerTest.DeleteServerTest(router)
	ServerTest.GetLeaderTest(router)
	ServerTest.GetAllServersTest(router)
	ServerTest.GetDeleteInstances(router)
	ServerTest.CancelDeleteInstance(router)
	ServerTest.GetDeleteInstances(router)
}

/**
注册-选举-获取leader领导的实例-拉入黑名单-拉出黑名单
*/
func test2(router *gin.Engine) {
	ServerTest.RegisterTest(router)
	ServerTest.GetAllServersTest(router)
	ServerTest.ElectionTest(router)
	ServerTest.GetLeaderTest(router)
	ServerTest.GetFollowServer(router)
	ServerTest.DeleteServerByColony(router)
	ServerTest.CancelDeleteInstance(router)
	ServerTest.GetDeleteInstances(router)
}

/**
注册-获取指定条件下服务器列表-拉入黑名单-拉出黑名单
*/
func test3(router *gin.Engine) {
	ServerTest.RegisterTest(router)
	ServerTest.GetAllServersTest(router)
	ServerTest.GetInstancesByQuery(router)
	ServerTest.DeleteServerTest(router)
	ServerTest.CancelDeleteInstance(router)
}

/**
注册-心跳-拉入黑名单-拉出黑名单
*/
func test4(router *gin.Engine) {
	ServerTest.RegisterTest(router)
	var num int
	for i := 0; i < 400; i++ {
		num++
		ServerTest.ServerBeatTest(router) // 服务心跳
		fmt.Println(num)
	}
	ServerTest.GetAllServersTest(router)
	ServerTest.DeleteServerTest(router)
	ServerTest.CancelDeleteInstance(router)
}
