package LeaderAlgorithm

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/util"
)

func CreateLeader(List *util.LinkList[entity.ServerModel]) entity.ServerModel {
	F := AlgorithmFactory()
	return F(List)
}

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
