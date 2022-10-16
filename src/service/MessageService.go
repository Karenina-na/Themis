package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/service/LeaderAlgorithm"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"strconv"
	"strings"
	"time"
)

// RegisterServer
// @Description: 注册服务器
// @param        S 服务器模型
// @return       B 是否成功
// @return       E 错误
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
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.RegisterChan <- *S
	}
	util.Loglevel(util.Debug, "RegisterServer", "注册服务-"+util.Strval(*S))
	return true, nil
}

// FlashHeartBeat
// @Description: 心跳
// @param        model 服务器模型
// @return       B     是否成功
// @return       E     错误
func FlashHeartBeat(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("FlashHeartBeat-service", util.Strval(r))
		}
	}()
	if config.ServerBeat.ServerModelBeatEnable {
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

// Election
// @Description: 选举
// @param        model 服务器模型
// @return       B     是否成功
// @return       E     错误
func Election(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Election-service", util.Strval(r))
		}
	}()
	ServerNum, err := GetServersNumber(model)
	if err != nil {
		return false, err
	}
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].IsEmpty() && ServerNum > 1 {
		Bean.RoutinePool.CreateWork(func() (E error) {
			time.Sleep(time.Duration(config.ElectionTimeOut) * time.Second)
			Bean.Leaders.LeaderModelsListRWLock.Lock()
			if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].IsEmpty() {
				Bean.Leaders.LeaderModelsListRWLock.Unlock()
				return nil
			}
			util.Loglevel(util.Info, "Election", "选举超时-"+model.Namespace+"."+model.Colony+"-选票数："+
				util.Strval(Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length()))
			Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Clear()
			Bean.Leaders.LeaderModelsListRWLock.Unlock()
			return nil
		}, func(Message error) {
			exception.HandleException(exception.NewUserError("Election-goroutine-service", Message.Error()))
		})
	}
	if !Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Contain(*model) {
		Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Append(*model)
	}
	if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length() <= ServerNum/2 {
		Bean.Leaders.LeaderModelsListRWLock.Unlock()
		return true, nil
	}
	util.Loglevel(util.Info, "Election", "选举开始-"+model.Namespace+"."+model.Colony+"-选票数："+
		util.Strval(Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length()))
	Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Clear()
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	Bean.Servers.ServerModelsListRWLock.RLock()
	ChoiceList := util.NewLinkList[entity.ServerModel]()
	for colonyMap, servers := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(colonyMap, "::")
		colony := str[0]
		if colony == model.Colony {
			servers.Iterator(func(index int, m entity.ServerModel) {
				ChoiceList.Append(m)
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	leader := LeaderAlgorithm.CreateLeader(ChoiceList)
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.LeaderChan <- leader
	}
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony] = leader
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "Election", "选举完成，发起通信-leader:"+leader.IP)
	ChoiceList.Iterator(func(index int, m entity.ServerModel) {
		util.Loglevel(util.Error, "Election", util.Strval(m))
		Bean.RoutinePool.CreateWork(func() (E error) {
			defer func() {
				r := recover()
				if r != nil {
					E = exception.NewSystemError("udp-message-goroutine", util.Strval(r))
				}
			}()
			//return m.SendMessageUDP(leader, config.Port.UDPPort, config.Port.UDPTimeOut)	//暂时注释

			port, _ := strconv.Atoi(m.Port)                                                //不要的
			return m.SendMessageUDP(leader, strconv.Itoa(port+10), config.Port.UDPTimeOut) //不要的

		}, func(err error) {
			exception.HandleException(exception.NewUserError("udp-message", " UDP通信错误 "+err.Error()+" "+util.Strval(model)))
		})
	})
	return true, nil
}

// GetLeader
// @Description: 获取leader
// @param        model 服务器模型
// @return       m     leader
// @return       E     错误
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

// GetServers
// @Description: 获取服务器列表
// @param        model 服务器模型
// @return       m     服务器列表
// @return       E     错误
func GetServers(serverModel *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServers-service", util.Strval(r))
		}
	}()
	list := make([]entity.ServerModel, 0, 100)
	Bean.Servers.ServerModelsListRWLock.RLock()
	for name, L := range Bean.Servers.ServerModelsList[serverModel.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == serverModel.Colony {
			L.Iterator(func(index int, model entity.ServerModel) {
				if serverModel.IP != model.IP {
					list = append(list, model)
				}
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return list, nil
}

// GetServersNumber
// @Description: 获取服务器数量
// @param        model 服务器模型
// @return       num   服务器数量
// @return       E     错误
func GetServersNumber(model *entity.ServerModel) (num int, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetServersNumber-service", util.Strval(r))
		}
	}()
	number := 0
	Bean.Servers.ServerModelsListRWLock.RLock()
	for name, List := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			number += List.Length()
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return number, nil
}
