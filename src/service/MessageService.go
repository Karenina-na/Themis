package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/entity/util"
	"Themis/src/service/LeaderAlgorithm"
	"encoding/json"
	"net"
)

func RegisterServer(S *entity.ServerModel) bool {
	if S.Namespace == "" {
		S.Namespace = "default"
	}
	if S.Colony == "" {
		S.Colony = "default"
	}
	ServerModelQueue <- *S
	return true
}

func FlashHeartBeat(model *entity.ServerModel) bool {
	ServerModelBeatQueueLock.Lock()
	ServerModelBeatQueue <- *model
	ServerModelBeatQueueLock.Unlock()
	return true
}

func Election(model *entity.ServerModel) bool {
	leader := LeaderAlgorithm.CreateLeader(ServerModelList[model.Namespace])
	Leader = leader
	List := ServerModelList[leader.Namespace]
	for _, list := range List {
		for i := 0; i < list.Length(); i++ {
			server := list.Get(i)
			RoutinePool.CreateWork(func() {
				udpAddr, _ := net.ResolveUDPAddr("udp", server.IP+":"+config.UDPPort)
				conn, err := net.DialUDP("udp", nil, udpAddr)
				if err != nil {
					util.Loglevel(util.Info, "Election-Service",
						"UDP通信错误 "+err.Error()+" "+util.Strval(server))
				} else {
					data, _ := json.Marshal(leader)
					_, _ = conn.Write(data)
				}
			})
		}
	}
	return true
}

func GetLeader() entity.ServerModel {
	return Leader
}

func GetServers(model *entity.ServerModel) []entity.ServerModel {
	list := make([]entity.ServerModel, 0, 100)
	ServerModelListRWLock.RLock()
	defer ServerModelListRWLock.RUnlock()
	for _, L := range ServerModelList[model.Namespace] {
		for i := 0; i < L.Length(); i++ {
			server := L.Get(i)
			if model.IP != server.IP {
				list = append(list, server)
			}
		}
	}
	return list
}
