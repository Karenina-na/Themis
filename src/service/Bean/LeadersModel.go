package Bean

import (
	"Themis/src/entity"
	"sync"
)

type LeadersModel struct {
	// LeaderModelsList 记账人
	LeaderModelsList map[string]map[string]entity.ServerModel
	// LeaderModelsListRWLock 记账人读写锁
	LeaderModelsListRWLock sync.RWMutex
}

func NewLeadersModel() *LeadersModel {
	return &LeadersModel{
		LeaderModelsList:       make(map[string]map[string]entity.ServerModel),
		LeaderModelsListRWLock: sync.RWMutex{},
	}
}
