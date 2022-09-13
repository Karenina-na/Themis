package LeaderAlgorithm

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/util"
)

// CreateLeader
// @Description: 生成Leader
// @param        List 服务器列表
// @return       entity.ServerModel Leader
func CreateLeader(List *util.LinkList[entity.ServerModel]) entity.ServerModel {
	F := AlgorithmFactory()
	return F(List)
}

// AlgorithmFactory
// @Description: 生成算法
// @return       func(List *util.LinkList[entity.ServerModel]) entity.ServerModel
func AlgorithmFactory() func(List *util.LinkList[entity.ServerModel]) entity.ServerModel {
	switch config.CreateLeaderAlgorithm {
	case "RandomAlgorithmCreateLeader":
		return RandomAlgorithmCreateLeader
	//case "xxx":
	//	return xxx
	default:
		return RandomAlgorithmCreateLeader
	}
}
