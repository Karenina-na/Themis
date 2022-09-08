package main

import (
	Factory "Themis/src/Factory"
	"Themis/src/util/encryption"
	"Themis/test/Base"
	"Themis/test/ServerTest"
	"crypto"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/goleak"
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
	defer goleak.VerifyNone(t)
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

	Factory.ThemisCloseFactory()

	/*哈希算法*/
	test5()
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

/*哈希算法*/
func test5() {
	//sha1
	sha1 := encryption.Sha1("123456")
	fmt.Println(sha1)
	//sha256
	sha256 := encryption.Sha256("123456")
	fmt.Println(sha256)
	//sha512
	sha512 := encryption.Sha512("123456")
	fmt.Println(sha512)
	//bcrypt
	bcrypt := encryption.Bcrypt("123456")
	fmt.Println(bcrypt)
	//base64
	base64 := encryption.Base64("base64")
	fmt.Println(base64)
	base64Decode := encryption.Base64Decode(base64)
	fmt.Println(base64Decode)
	//ASE
	var PwdKey = []byte("DIS**#KKKDJJSKDI")
	ASE := encryption.AESEncrypt("ASE", PwdKey)
	fmt.Println(ASE)
	ASEDecrypt := encryption.AESDecrypt(ASE, PwdKey)
	fmt.Println(ASEDecrypt)
	//RSA
	privateKey, publicKey := encryption.RSACreateKey()
	RSA := encryption.RSAEncrypt(publicKey, "RSA", crypto.SHA256.New())
	fmt.Println(RSA)
	RSADecrypt := encryption.RSADecrypt(privateKey, RSA, crypto.SHA256)
	fmt.Println(RSADecrypt)
}
