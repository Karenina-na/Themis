package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/LeaderAlgorithm"
	"Themis/src/util"
	"encoding/json"
	"net"
	"strings"
)

func RegisterServer(S *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("RegisterServer-service", util.Strval(r))
		}
	}()
	if S.Namespace == "" {
		S.Namespace = "default"
	}
	if S.Colony == "" {
		S.Colony = "default"
	}
	ServerModelQueue <- *S
	util.Loglevel(util.Debug, "RegisterServer", "注册服务-"+util.Strval(*S))
	return true, nil
}

func FlashHeartBeat(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("FlashHeartBeat-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "FlashHeartBeat", "刷新心跳-"+util.Strval(*model))
	ServerModelBeatQueueLock.Lock()
	ServerModelBeatQueue <- *model
	ServerModelBeatQueueLock.Unlock()
	return true, nil
}

func Election(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Election-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Election", "选举开始")
	leader := LeaderAlgorithm.CreateLeader(ServerModelList[model.Namespace])
	Leaders[model.Namespace] = leader
	List := ServerModelList[leader.Namespace]
	util.Loglevel(util.Debug, "Election", "选举完成，发起通信-leader:"+leader.IP)
	for _, list := range List {
		for i := 0; i < list.Length(); i++ {
			server := list.Get(i)
			RoutinePool.CreateWork(func() (E error) {
				defer func() {
					r := recover()
					if r != nil {
						E = exception.NewSystemError("udp-message-goroutine", util.Strval(r))
					}
				}()
				udpAddr, _ := net.ResolveUDPAddr("udp", server.IP+":"+config.UDPPort)
				conn, err := net.DialUDP("udp", nil, udpAddr)
				if err != nil {
					return exception.NewUserError("udp-message", "UDP通信错误"+err.Error()+util.Strval(server))
				} else {
					data, _ := json.Marshal(leader)
					_, _ = conn.Write(data)
				}
				return nil
			}, func(Message error) {
				exception.HandleException(Message)
			})
		}
	}
	return true, nil
}

func GetLeader(model *entity.ServerModel) (m entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetLeader-service", util.Strval(r))
		}
	}()
	return Leaders[model.Namespace], nil
}

func GetServers(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServers-service", util.Strval(r))
		}
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

func GetServersNumber(model *entity.ServerModel) (num int, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServersNumber-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.RLock()
	defer ServerModelListRWLock.RUnlock()
	for name, List := range ServerModelList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			return List.Length(), nil
		}
	}
	return 0, nil
}
