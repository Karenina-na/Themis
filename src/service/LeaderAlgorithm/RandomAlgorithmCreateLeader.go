package LeaderAlgorithm

import (
	"Themis/src/entity"
	"Themis/src/util"
	"math/rand"
	"time"
)

func RandomAlgorithmCreateLeader(List *util.LinkList[entity.ServerModel]) entity.ServerModel {
	rand.Seed(time.Now().Unix())
	num := rand.Int() % List.Length()
	server := entity.NewServerModel()
	for i := 0; ; i++ {
		if num == 0 {
			server.Namespace = List.Get(i).Namespace
			server.Colony = List.Get(i).Colony
			server.IP = List.Get(i).IP
			server.Port = List.Get(i).Port
			server.Name = List.Get(i).Name
			server.Time = List.Get(i).Time
			break
		}
		num--
	}
	return *server
}
