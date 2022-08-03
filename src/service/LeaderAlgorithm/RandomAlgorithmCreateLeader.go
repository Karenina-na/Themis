package LeaderAlgorithm

import (
	"Themis/src/entity"
	"Themis/src/entity/util"
	"math/rand"
	"time"
)

func RandomAlgorithmCreateLeader(List map[string]*util.LinkList[entity.ServerModel]) entity.ServerModel {
	rand.Seed(time.Now().Unix())
	num := rand.Int() % len(List)
	var servers *util.LinkList[entity.ServerModel]
	for _, v := range List {
		if num == 0 {
			servers = v
			break
		}
		num--
	}
	return servers.Get(rand.Int() % servers.Length())
}
