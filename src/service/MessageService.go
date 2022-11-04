package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	Factory "Themis/src/pool"
	"Themis/src/service/Bean"
	"Themis/src/service/LeaderAlgorithm"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
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

	//添加实例信息，注册进入队列
	Bean.InstanceList.Append(*S)
	Bean.ServersQueue.Enqueue(S)

	//如果开启数据同步，则注册同步服务
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.RegisterChan.Enqueue(S)
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

		//刷新心跳服务
		for flag == false {
			Bean.ServersBeatQueue.Operate(func() {

				//队列不满则直接添加
				if !Bean.ServersBeatQueue.IsFull() {
					Bean.ServersBeatQueue.Enqueue(model)
					flag = true
				}
			})
			time.Sleep(time.Millisecond)
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
	//获取服务数量
	ServerNum, err := GetServersNumber(model)
	if err != nil {
		return false, err
	}

	//如果服务数量为1，且为第一次发起选举，创建管理选举超时线程
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].IsEmpty() && ServerNum > 1 {
		Factory.RoutinePool.CreateWork(func() (E error) {
			time.Sleep(time.Duration(config.ElectionTimeOut) * time.Second)
			Bean.Leaders.LeaderModelsListRWLock.Lock()
			if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].IsEmpty() {
				Bean.Leaders.LeaderModelsListRWLock.Unlock()
				return nil
			}
			util.Loglevel(util.Info, "Election", "选举超时-"+model.Namespace+"."+model.Colony+"-选票数："+
				util.Strval(Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length()))

			//选举超时，清空选票
			Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Clear()
			Bean.Leaders.LeaderModelsListRWLock.Unlock()
			return nil
		}, func(Message error) {
			exception.HandleException(exception.NewUserError("Election-goroutine-service", Message.Error()))
		})
	}

	//添加选票
	if !Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Contain(*model) {
		Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Append(*model)
	}

	//如果选票数小于等于服务数量的一半，则继续等待
	if Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length() <= ServerNum/2 {
		Bean.Leaders.LeaderModelsListRWLock.Unlock()
		return true, nil
	}

	//如果选票数大于服务数量的一半，则开启选举
	util.Loglevel(util.Info, "Election", "选举开始-"+model.Namespace+"."+model.Colony+"-选票数："+
		util.Strval(Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Length()))

	//选举开始，清空选票
	Bean.Leaders.ElectionServers[model.Namespace][model.Colony].Clear()
	Bean.Leaders.LeaderModelsListRWLock.Unlock()

	//获取参与选举的服务
	Bean.Servers.ServerModelsListRWLock.RLock()
	ChoiceList := util.NewLinkList[entity.ServerModel]()
	for colonyMap, servers := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(colonyMap, "::")
		colony := str[0]
		if colony == model.Colony {
			servers.Iterator(func(index int, m entity.ServerModel) {
				ChoiceList.Append(*m.Clone())
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()

	//选举算法，得到领导者
	leader := LeaderAlgorithm.CreateLeader(ChoiceList)

	//集群同步
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.LeaderChan.Enqueue(&leader)
	}

	//设置领导者
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony] = leader
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "Election", "选举完成，发起通信-leader:"+leader.IP)

	//发起广播通信
	ChoiceList.Iterator(func(index int, m entity.ServerModel) {
		Factory.RoutinePool.CreateWork(func() (E error) {
			defer func() {
				r := recover()
				if r != nil {
					E = exception.NewSystemError("udp-message-goroutine", util.Strval(r))
				}
			}()
			return m.SendMessageUDP(leader, config.Port.UDPTimeOut)

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
// @Description: 获取领导者领导的服务器列表
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
			//迭代获取与领导者在相同集群的服务器
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
		//迭代获取与领导者在相同集群的服务器数量
		if colony == model.Colony {
			number += List.Length()
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return number, nil
}
