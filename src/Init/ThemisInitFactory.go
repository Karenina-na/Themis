package FactoryInit

import (
	"Themis/docs"
	"Themis/src/config"
	"Themis/src/entity/util"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/service"
)

func ThemisInitFactory() (E any) {
	util.LoggerInit(func(r any) {
		exception.HandleException(exception.NewServicePanic("日志模块", util.Strval(r)))
	})
	if err := config.InitConfig(); err != nil {
		return err
	}
	if err := config.SwaggerConfig(docs.SwaggerInfo); err != nil {
		return err
	}
	if config.DatabaseEnable {
		if err := mapper.FactoryInit(); err != nil {
			return err
		}
	}
	if err := service.ServerInitFactory(); err != nil {
		return err
	}
	return nil
}
