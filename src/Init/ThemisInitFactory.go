package FactoryInit

import (
	"Themis/docs"
	"Themis/src/config"
	"Themis/src/mapper"
	"Themis/src/service"
)

func ThemisInitFactory() {
	config.SwaggerConfig(docs.SwaggerInfo)
	if config.DatabaseEnable {
		mapper.FactoryInit()
	}
	service.ServerInitFactory()
}
