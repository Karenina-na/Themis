package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/service/LeaderAlgorithm"
	"Themis/src/util"
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
	Bean.ServersQueue <- *S
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
	if config.ServerModelBeatEnable {
		util.Loglevel(util.Debug, "FlashHeartBeat", "刷新心跳-"+util.Strval(*model))
		flag := false
		for flag == false {
			select {
			case Bean.ServersBeatQueue <- *model:
				flag = true
				break
			default:
				time.Sleep(time.Millisecond)
			}
		}
		return true, nil
	}
	util.Loglevel(util.Debug, "FlashHeartBeat", "心跳服务未开启")
	return false, nil
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
	Bean.Servers.ServerModelsListRWLock.RLock()
	ChoiceList := util.NewLinkList[entity.ServerModel]()
	for colonyMap, servers := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(colonyMap, "::")
		colony := str[0]
		if colony == model.Colony {
			servers.Iterator(func(index int, model entity.ServerModel) {
				ChoiceList.Append(model)
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	leader := LeaderAlgorithm.CreateLeader(ChoiceList)
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony] = leader
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "Election", "选举完成，发起通信-leader:"+leader.IP)
	ChoiceList.Iterator(func(index int, model entity.ServerModel) {
		Bean.RoutinePool.CreateWork(func() (E error) {
			defer func() {
				r := recover()
				if r != nil {
					E = exception.NewSystemError("udp-message-goroutine", util.Strval(r))
				}
			}()
			E = model.SendMessageUDP(leader, config.UDPPort, config.UDPTimeOut)
			return E
		}, func(err error) {
			exception.HandleException(exception.NewUserError("udp-message", " UDP通信错误 "+err.Error()+" "+util.Strval(model)))
		})
	})
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
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	leader := Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony]
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
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
	Bean.Servers.ServerModelsListRWLock.RLock()
	for _, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		L.Iterator(func(index int, model entity.ServerModel) {
			if model.IP != model.IP {
				list = append(list, model)
			}
		})
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
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
	Bean.Servers.ServerModelsListRWLock.RLock()
	for name, List := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			return List.Length(), nil
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return 0, nil
}
