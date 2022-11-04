package Init

import (
	"Themis/docs"
	"Themis/src/Factory/image"
	"Themis/src/config"
	"Themis/src/config/LoadConfigFactory"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/pool"
	"Themis/src/service"
	"Themis/src/sync"
	"Themis/src/util"
	"time"
)

// ThemisInitFactory
// @Description: Themis初始化工厂
// @param        arg 初始化级别
func ThemisInitFactory(arg *string) {
	//初始化日志模块
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

	//初始化配置模块
	if err := LoadConfigFactory.InitConfig(); err != nil {
		exception.HandleException(err)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化配置文件")

	//初始化协程池
	if err := pool.InitRoutinePool(); err != nil {
		exception.HandleException(err)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化协程池")

	//初始化文档模块
	if err := config.SwaggerConfig(docs.SwaggerInfo); err != nil {
		exception.HandleException(err)
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化Swagger文档")

	//初始化数据库模块
	if config.Persistence.PersistenceEnable {
		util.Loglevel(util.Debug, "ThemisInitFactory", "初始化数据库")
		if err := mapper.InitMapper(); err != nil {
			exception.HandleException(err)
		}
	}

	//初始化集群模块
	if config.Cluster.ClusterEnable {
		util.Loglevel(util.Debug, "ThemisInitFactory", "初始化集群")
		if err := sync.InitSync(); err != nil {
			exception.HandleException(err)
		}
	}
	util.Loglevel(util.Debug, "ThemisInitFactory", "初始化服务")

	//初始化服务模块
	if err := service.InitServer(); err != nil {
		exception.HandleException(err)
	}

	//画图
	time.Sleep(time.Second * 1)
	image.Factory()
	time.Sleep(time.Second * 1)
}
