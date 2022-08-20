package FactoryInit

import (
	"Themis/docs"
	"Themis/src/config"
	"Themis/src/mapper"
	"Themis/src/service"
)

func ThemisInitFactory() (E any) {
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
