package sync

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"time"
)

//
// Leader
// @Description: 领导者
// @return       E 异常
//
func Leader() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Leader-sync", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Leader-sync", "Leader状态")
	syncBean.Status = syncBean.LEADER
	syncBean.Leader.SetLeaderModel(syncBean.Name, config.Cluster.IP, config.Cluster.Port, config.Port.CenterPort)
	SyncRoutineBool := make(chan bool, 0)
	Bean.RoutinePool.CreateWork(func() (E error) {
		util.Loglevel(util.Debug, "Leader-sync", "leader数据发送协程启动")
		for {
			select {
			case <-time.After(time.Millisecond * time.Duration(config.Cluster.LeaderSyncTime)):
				instances, deleteInstances, leaderInstances, e := CreateSendSyncData()
				if e != nil {
					return e
				}
				m := syncBean.NewMessageModel()
				syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
					m.SetMessageMode(syncBean.Term, syncBean.Status,
						instances, deleteInstances, leaderInstances,
						config.Port.CenterPort,
						config.Cluster.IP, config.Cluster.Port, value.IP, value.Port, false)
					syncBean.UdpSendMessage <- *m
					if config.Cluster.TrackEnable {
						util.Loglevel(util.Debug, "Leader-sync", "leader数据发送协程发送数据"+util.Strval(m.UDPTargetAddress))
					}
				})
			case <-SyncRoutineBool:
				util.Loglevel(util.Info, "Leader-sync", "leader数据发送协程关闭")
				return nil
			case <-syncBean.CloseChan:
				util.Loglevel(util.Info, "Leader-sync", "leader数据发送协程关闭")
				return nil
			}
		}
	}, func(Message error) {
		exception.HandleException(Message)
	})
	for {
		select {
		case m := <-syncBean.UdpReceiveMessage:
			if m.Term >= syncBean.Term {
				switch m.Status {
				case syncBean.LEADER:
					if m.Term > syncBean.Term {
						syncBean.Term = m.Term
						syncBean.Leader.SetLeaderModel(m.Name, m.UDPAddress.IP, m.UDPAddress.Port, m.ServicePort)
						close(SyncRoutineBool)
						util.Loglevel(util.Info, "Leader-sync", "Leader卸任,收到其他leader，成为FOLLOW")
						E := ChangeToFollow()
						if E != nil {
							return E
						}
						return nil
					}
					if m.Term == syncBean.Term {
						close(SyncRoutineBool)
						util.Loglevel(util.Info, "Leader-sync", "Leader卸任,收到其他相同任期的leader，成为FOLLOW")
						E := ChangeToFollow()
						if E != nil {
							return E
						}
						return nil
					}
				case syncBean.CANDIDATE:
					if m.Term > syncBean.Term {
						message := syncBean.NewMessageModel()
						message.SetMessageMode(syncBean.Term, syncBean.Status,
							nil, nil, nil,
							config.Port.CenterPort,
							config.Cluster.IP, config.Cluster.Port,
							m.UDPAddress.IP, m.UDPAddress.Port, true)
						syncBean.UdpSendMessage <- *message
						syncBean.Term = m.Term
						close(SyncRoutineBool)
						util.Loglevel(util.Info, "Leader-sync", "Leader卸任,收到更高任期的CANDIDATE，成为FOLLOW")
						E := ChangeToFollow()
						if E != nil {
							return E
						}
						return nil
					}
				case syncBean.FOLLOW:
					if m.Term > syncBean.Term {
						syncBean.Term = m.Term
						close(SyncRoutineBool)
						util.Loglevel(util.Info, "Leader-sync", "Leader卸任,收到更高任期的FOLLOW，成为FOLLOW")
						E := ChangeToFollow()
						if E != nil {
							return E
						}
						return nil
					}
				}
			}
		case <-syncBean.CloseChan:
			close(SyncRoutineBool)
			util.Loglevel(util.Debug, "Leader-sync", "Leader协程关闭")
			return nil
		}
	}
}

//
// CreateSendSyncData
// @Description: 创建发送同步数据
// @return       instances  实例列表
// @return       list       删除实例列表
// @return       leaderList leader实例列表
// @return       E          错误
//
func CreateSendSyncData() (instances []entity.ServerModel,
	list []entity.ServerModel, leaderList []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSendSyncData-sync", util.Strval(r))
		}
	}()
	LeaderList := make([]entity.ServerModel, 0)
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	for _, Namespace := range Bean.Leaders.LeaderModelsList {
		for _, Leader := range Namespace {
			LeaderList = append(LeaderList, Leader)
		}
	}
	Instances := make([]entity.ServerModel, 0)
	Bean.InstanceList.Iterator(func(index int, value entity.ServerModel) {
		Instances = append(Instances, value)
	})
	DeleteInstancesList := make([]entity.ServerModel, 0)
	Bean.DeleteInstanceList.Iterator(func(index int, value entity.ServerModel) {
		DeleteInstancesList = append(DeleteInstancesList, value)
	})
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
	return Instances, DeleteInstancesList, LeaderList, nil
}
