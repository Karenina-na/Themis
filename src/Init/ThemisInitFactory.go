package FactoryInit

import (
	"Themis/docs"
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/service"
	"Themis/src/util"
)

func ThemisInitFactory() {
	util.LoggerInit(func(r any) {
		exception.HandleException(exception.NewSystemError("日志模块", util.Strval(r)))
	}, util.Debug)
	if err := config.InitConfig(); err != nil {
		exception.HandleException(err)
	}
	if err := config.SwaggerConfig(docs.SwaggerInfo); err != nil {
		exception.HandleException(err)
	}
	if config.DatabaseEnable {
		if err := mapper.InitMapper(); err != nil {
			exception.HandleException(err)
		}
	}
	if err := service.InitServer(); err != nil {
		exception.HandleException(err)
	}
}
