package Init

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/service"
	"Themis/src/sync"
)

func ThemisCloseFactory() {
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
