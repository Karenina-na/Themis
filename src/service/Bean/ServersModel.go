package Bean

import (
	"Themis/src/entity"
	"Themis/src/util"
	"sync"
)

type ServersModel struct {
	// ServerModelsList 服务模型
	ServerModelsList map[string]map[string]*util.LinkList[entity.ServerModel]
	// ServerModelsListRWLock 服务模型读写锁
	ServerModelsListRWLock sync.RWMutex
}

func NewServersModel() *ServersModel {
	return &ServersModel{
		ServerModelsList:       make(map[string]map[string]*util.LinkList[entity.ServerModel]),
		ServerModelsListRWLock: sync.RWMutex{},
	}
}
