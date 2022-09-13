package Base

import (
	"Themis/src/entity"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

// Get
// @Description: get请求
// @param        uri    请求地址
// @param        router 路由
// @return       *httptest.ResponseRecorder 响应
func Get(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

// Post
// @Description: post请求
// @param        uri    请求地址
// @param        param  请求参数
// @param        router 路由
// @return       *httptest.ResponseRecorder 响应
func Post(uri string, param entity.ServerModel, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest("POST", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Content-Type", "application/json")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}

// Put
// @Description: put请求
// @param        uri    请求地址
// @param        param  请求参数
// @param        router 路由
// @return       *httptest.ResponseRecorder 响应
func Put(uri string, param entity.ServerModel, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest("PUT", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Content-Type", "application/json")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}

// Delete
// @Description: delete请求
// @param        uri    请求地址
// @param        param  请求参数
// @param        router 路由
// @return       *httptest.ResponseRecorder 响应
func Delete(uri string, param entity.ServerModel, router *gin.Engine) *httptest.ResponseRecorder {
	jsonByte, _ := json.Marshal(param)
	req := httptest.NewRequest("DELETE", uri, bytes.NewReader(jsonByte))
	req.Header.Add("Content-Type", "application/json")
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应handler接口
	router.ServeHTTP(w, req)
	return w
}
