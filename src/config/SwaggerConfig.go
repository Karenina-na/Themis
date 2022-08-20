package config

import "github.com/swaggo/swag"

func SwaggerConfig(docs *swag.Spec) (E any) {
	defer func() {
		E = recover()
	}()
	docs.Title = "Themis API"
	docs.Description = "分布式记账系统调度中心"
	docs.Version = "1.0"
	docs.Host = Port
	docs.Schemes = []string{"http", "https"}
	docs.BasePath = "/api/v1"
	return nil
}
