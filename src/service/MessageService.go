package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/entity/util"
	"Themis/src/exception"
	"Themis/src/service/LeaderAlgorithm"
	"encoding/json"
	"net"
)

func RegisterServer(S *entity.ServerModel) (B bool, E any) {
	defer func() {
		E = recover()
	}()
	if S.Namespace == "" {
		S.Namespace = "default"
	}
	if S.Colony == "" {
		S.Colony = "default"
	}
	ServerModelQueue <- *S
	return true, nil
}

func FlashHeartBeat(model *entity.ServerModel) (B bool, E any) {
	defer func() {
		E = recover()
	}()
	ServerModelBeatQueueLock.Lock()
	ServerModelBeatQueue <- *model
	ServerModelBeatQueueLock.Unlock()
	return true, nil
}

func Election(model *entity.ServerModel) (B bool, E any) {
	defer func() {
		E = recover()
	}()
	leader := LeaderAlgorithm.CreateLeader(ServerModelList[model.Namespace])
	Leader = leader
	List := ServerModelList[leader.Namespace]
	for _, list := range List {
		for i := 0; i < list.Length(); i++ {
			server := list.Get(i)
			RoutinePool.CreateWork(func() (E any) {
				defer func() {
					E = recover()
				}()
				udpAddr, _ := net.ResolveUDPAddr("udp", server.IP+":"+config.UDPPort)
				conn, err := net.DialUDP("udp", nil, udpAddr)
				if err != nil {
					panic(exception.NewServicePanic("Election", "UDP通信错误"+err.Error()+util.Strval(server)))
				} else {
					data, _ := json.Marshal(leader)
					_, _ = conn.Write(data)
				}
				return nil
			}, func(Message any) {
				panic(exception.NewServicePanic("Election-Service", "goroutine错误"+util.Strval(Message)))
			})
		}
	}
	return true, nil
}

func GetLeader() (m entity.ServerModel, E any) {
	defer func() {
		E = recover()
	}()
	return Leader, nil
}

func GetServers(model *entity.ServerModel) (m []entity.ServerModel, E any) {
	defer func() {
		E = recover()
	}()
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
	return list, nil
}
