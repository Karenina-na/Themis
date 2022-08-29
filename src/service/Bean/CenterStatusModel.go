package Bean

import (
	"Themis/src/entity"
	"sync"
)

type CenterStatusModel struct {
	// CenterStatusInfo 服务器状态
	CenterStatusInfo *entity.ComputerInfoModel
	// CenterStatusInfoLock  服务器状态读写锁
	CenterStatusInfoLock sync.RWMutex
}

func NewCenterStatusModel() *CenterStatusModel {
	return &CenterStatusModel{
		CenterStatusInfo:     entity.NewComputerInfoModel(),
		CenterStatusInfoLock: sync.RWMutex{},
	}
}
