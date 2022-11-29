package Init

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/pool"
	"Themis/src/service"
	"Themis/src/sync"
)

// ThemisCloseFactory
// @Description: ThemisCloseFactory
func ThemisCloseFactory() {
	//关闭数据库模块
	if config.Persistence.PersistenceEnable {
		err := mapper.CloseMapper()
		if err != nil {
			exception.HandleException(err)
		}
	}

	//关闭集群模块
	if config.Cluster.ClusterEnable {
		err := sync.Close()
		if err != nil {
			exception.HandleException(err)
		}
	}

	//关闭服务模块
	if err := service.Close(); err != nil {
		exception.HandleException(err)
	}

	//关闭协程池
	if err := pool.CloseRoutinePool(); err != nil {
		exception.HandleException(err)
	}
}
