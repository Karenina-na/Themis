package Init

import (
	"Themis/docs"
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/service"
	"Themis/src/sync"
	"Themis/src/util"
)

func ThemisInitFactory(arg *string) {
	if *arg == "debug" {
		util.LoggerInit(func(r any) {
			exception.HandleException(exception.NewSystemError("日志模块", util.Strval(r)))
		}, util.Debug)
	} else {
		util.LoggerInit(func(r any) {
			exception.HandleException(exception.NewSystemError("日志模块", util.Strval(r)))
		}, util.Info)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化日志模块")
	if err := config.InitConfig(); err != nil {
		exception.HandleException(err)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化配置文件")
	if err := config.SwaggerConfig(docs.SwaggerInfo); err != nil {
		exception.HandleException(err)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化Swagger文档")
	if config.Persistence.PersistenceEnable {
		util.Loglevel(util.Debug, "ThemisInitFactory", "初始化数据库")
		if err := mapper.InitMapper(); err != nil {
			exception.HandleException(err)
		}
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化服务")
	if err := service.InitServer(); err != nil {
		exception.HandleException(err)
	}
	if config.Cluster.ClusterEnable {
		util.Loglevel(util.Debug, "ThemisInitFactory", "初始化集群")
		if err := sync.InitSync(); err != nil {
			exception.HandleException(err)
		}
	}
}
