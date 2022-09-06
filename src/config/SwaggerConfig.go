package config

import (
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/swaggo/swag"
)

//
// SwaggerConfig
// @Description: Swagger配置
// @param        docs 文档目录
// @return       E    错误
//
func SwaggerConfig(docs *swag.Spec) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SwaggerConfig-config", util.Strval(r))
		}
	}()
	docs.Title = "Themis API"
	docs.Description = "分布式记账系统调度中心"
	docs.Version = "1.0"
	docs.Host = "localhost:" + Port.CenterPort
	docs.Schemes = []string{"http", "https"}
	docs.BasePath = "/api/v1"
	return nil
}
