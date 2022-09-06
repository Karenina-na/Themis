package Bean

import (
	"Themis/src/entity"
	"sync"
)

// CenterStatusModel is the model of center status
type CenterStatusModel struct {
	// CenterStatusInfo 服务器状态
	CenterStatusInfo *entity.ComputerInfoModel
	// CenterStatusInfoLock  服务器状态读写锁
	CenterStatusInfoLock sync.RWMutex
}

//
// NewCenterStatusModel
// @Description: 生成CenterStatusModel
// @return       *CenterStatusModel 返回CenterStatusModel
//
func NewCenterStatusModel() *CenterStatusModel {
	return &CenterStatusModel{
		CenterStatusInfo:     entity.NewComputerInfoModel(),
		CenterStatusInfoLock: sync.RWMutex{},
	}
}
