package Base

import (
	"Themis/src/entity"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func Get(uri string, router *gin.Engine) *httptest.ResponseRecorder {
	// 构造get请求
	req := httptest.NewRequest("GET", uri, nil)
	// 初始化响应
	w := httptest.NewRecorder()
	// 调用相应的handler接口
	router.ServeHTTP(w, req)
	return w
}

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
