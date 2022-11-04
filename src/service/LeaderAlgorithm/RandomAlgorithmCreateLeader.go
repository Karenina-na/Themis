package LeaderAlgorithm

import (
	"Themis/src/entity"
	"Themis/src/util"
	"math/rand"
	"time"
)

// RandomAlgorithmCreateLeader
// @Description: 随机选举算法
// @param        List 节点列表
// @return       entity.ServerModel 选举出的Leader节点
func RandomAlgorithmCreateLeader(List *util.LinkList[entity.ServerModel]) entity.ServerModel {
	rand.Seed(time.Now().Unix())
	num := rand.Int() % List.Length()
	return List.Get(num)
}
