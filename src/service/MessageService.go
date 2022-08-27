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
	"time"
)

// RegisterServer 注册服务器
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

// FlashHeartBeat 心跳
func FlashHeartBeat(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("FlashHeartBeat-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "FlashHeartBeat", "刷新心跳-"+util.Strval(*model))
	flag := false
	for flag == false {
		select {
		case ServerModelBeatQueue <- *model:
			flag = true
			break
		default:
			time.Sleep(time.Millisecond)
		}
	}
	return true, nil
}

// Election 发起选举
func Election(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Election-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Election", "选举开始")
	ServerModelListRWLock.RLock()
	ChoiceList := util.NewLinkList[entity.ServerModel]()
	for colonyMap, servers := range ServerModelList[model.Namespace] {
		str := strings.Split(colonyMap, "::")
		colony := str[0]
		if colony == model.Colony {
			for i := 0; i < servers.Length(); i++ {
				ChoiceList.Append(servers.Get(i))
			}
		}
	}
	ServerModelListRWLock.RUnlock()
	leader := LeaderAlgorithm.CreateLeader(ChoiceList)
	LeadersRWLock.Lock()
	Leaders[model.Namespace][model.Colony] = leader
	LeadersRWLock.Unlock()
	util.Loglevel(util.Debug, "Election", "选举完成，发起通信-leader:"+leader.IP)
	for i := 0; i < ChoiceList.Length(); i++ {
		server := ChoiceList.Get(i)
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
				return exception.NewUserError("udp-message", " UDP通信错误 "+err.Error()+util.Strval(server))
			} else {
				data, _ := json.Marshal(leader)
				_, _ = conn.Write(data)
			}
			return nil
		}, func(Message error) {
			exception.HandleException(Message)
		})
	}
	return true, nil
}

// GetLeader 获取leader
func GetLeader(model *entity.ServerModel) (m entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetLeader-service", util.Strval(r))
		}
	}()
	LeadersRWLock.RLock()
	leader := Leaders[model.Namespace][model.Colony]
	LeadersRWLock.RUnlock()
	return leader, nil
}

// GetServers 获取leader领导的服务器列表
func GetServers(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServers-service", util.Strval(r))
		}
	}()
	list := make([]entity.ServerModel, 0, 100)
	ServerModelListRWLock.RLock()
	for _, L := range ServerModelList[model.Namespace] {
		for i := 0; i < L.Length(); i++ {
			server := L.Get(i)
			if model.IP != server.IP {
				list = append(list, server)
			}
		}
	}
	ServerModelListRWLock.RUnlock()
	return list, nil
}

// GetServersNumber 获取leader领导的服务器列表数量
func GetServersNumber(model *entity.ServerModel) (num int, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServersNumber-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.RLock()
	for name, List := range ServerModelList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			return List.Length(), nil
		}
	}
	ServerModelListRWLock.RUnlock()
	return 0, nil
}
