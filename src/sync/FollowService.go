package sync

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"math/rand"
	"time"
)

//
// Follow
// @Description: 跟随者
// @return       E error
//
func Follow() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Follow-sync", util.Strval(r))
		}
	}()
	syncBean.Status = syncBean.FOLLOW
	util.Loglevel(util.Debug, "Follow-sync", "FOLLOW状态")
	for {
		select {
		case <-time.After(time.Second * time.Duration(rand.Int()%
			int(config.Cluster.MaxFollowTimeOut-config.Cluster.MinFollowTimeOut)+int(config.Cluster.MinFollowTimeOut))):
			util.Loglevel(util.Debug, "Follow-sync", "FOLLOW超时，成为CANDIDATE")
			E := ChangeToCandidate()
			if E != nil {
				return E
			}
			return nil
		case m := <-syncBean.UdpReceiveMessage:
			if err := StatusOperatorFollow(&m); err != nil {
				return err
			}
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "Follow-sync", "FOLLOW协程退出")
			return nil
		}
	}
}

//
// StatusOperatorFollow
// @Description: FOLLOW状态下的操作
// @param        m *syncBean.MessageModel
// @return       E error
//
func StatusOperatorFollow(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("StatusOperatorFollow-sync", util.Strval(r))
		}
	}()
	switch m.Status {
	case syncBean.LEADER:
		if m.Term >= syncBean.Term {
			if config.Cluster.TrackEnable {
				util.Loglevel(util.Debug, "StatusOperatorFollow-sync", "收到LEADER信息-"+util.Strval(m.UDPAddress))
			}
			if m.UDPAddress.IP != syncBean.Leader.LeaderAddress.IP ||
				m.UDPAddress.Port != syncBean.Leader.LeaderAddress.Port {
				syncBean.Leader.SetLeaderModel(m.Name, m.UDPAddress.IP, m.UDPAddress.Port, m.ServicePort)
				syncBean.Term = m.Term
			}
			if err := CreateSyncRoutine(m); err != nil {
				return err
			}
		}
	case syncBean.CANDIDATE:
		if m.Term > syncBean.Term {
			if config.Cluster.TrackEnable {
				util.Loglevel(util.Debug, "StatusOperatorFollow-sync", "收到CANDIDATE信息-"+util.Strval(m.UDPAddress))
			}
			message := syncBean.NewMessageModel()
			message.SetMessageMode(syncBean.Term, syncBean.Status,
				nil, nil, nil,
				config.Port.CenterPort,
				config.Cluster.IP, config.Cluster.Port,
				m.UDPAddress.IP, m.UDPAddress.Port, true)
			if config.Cluster.TrackEnable {
				util.Loglevel(util.Debug, "StatusOperatorFollow-sync", "投票true")
			}
			syncBean.UdpSendMessage <- *message
		}
	}
	return nil
}

//
// CreateSyncRoutine
// @Description: 创建同步协程
// @param        m *syncBean.MessageModel
// @return       E error
//
func CreateSyncRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncRoutine-sync", util.Strval(r))
		}
	}()
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := DataSyncInstances(m.Message.Instances)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := DataSyncDelete(m.Message.DeleteInstances)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	Bean.RoutinePool.CreateWork(func() (E error) {
		err := DataSyncLeader(m.Message.Leaders)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}
