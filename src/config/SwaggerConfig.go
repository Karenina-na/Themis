package config

import (
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/swaggo/swag"
)

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
