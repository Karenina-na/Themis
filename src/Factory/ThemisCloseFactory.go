package Init

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/service"
	"Themis/src/sync"
)

func ThemisCloseFactory() {
	if config.Persistence.PersistenceEnable {
		err := mapper.CloseMapper()
		if err != nil {
			exception.HandleException(err)
		}
	}
	if config.Cluster.ClusterEnable {
		err := sync.Close()
		if err != nil {
			exception.HandleException(err)
		}
	}
	err := service.Close()
	if err != nil {
		exception.HandleException(err)
	}

}
