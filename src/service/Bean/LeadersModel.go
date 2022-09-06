package Bean

import (
	"Themis/src/entity"
	"sync"
)

// LeadersModel is a struct that contains a map of LeaderModels
type LeadersModel struct {
	// LeaderModelsList 记账人
	LeaderModelsList map[string]map[string]entity.ServerModel
	// LeaderModelsListRWLock 记账人读写锁
	LeaderModelsListRWLock sync.RWMutex
}

//
// NewLeadersModel
// @Description: 创建一个新的记账人模型
// @return       *LeadersModel 返回一个新的记账人模型
//
func NewLeadersModel() *LeadersModel {
	return &LeadersModel{
		LeaderModelsList:       make(map[string]map[string]entity.ServerModel),
		LeaderModelsListRWLock: sync.RWMutex{},
	}
}
