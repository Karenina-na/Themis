package leader

import (
	"Themis/src/config"
	"Themis/src/exception"
	Factory "Themis/src/pool"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"time"
)

// Leader
// @Description: 领导者
// @return       E 异常
func Leader() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Leader-leader", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Leader-leader", "Leader状态")
	syncBean.Status = syncBean.LEADER
	syncBean.Leader.SetLeaderModel(syncBean.Name, config.Cluster.IP, config.Cluster.Port, config.Port.CenterPort)

	//同步协程关闭标志
	SyncRoutineBool := make(chan bool, 0)

	//创建发送全部数据snapshot协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		util.Loglevel(util.Debug, "Leader-leader", "leader-snapshot数据发送协程启动")
		if err := SendDataGoroutineSnapshot(SyncRoutineBool); err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})

	//创建发送单个数据包协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		util.Loglevel(util.Debug, "Leader-leader", "leader-数据发送协程启动")
		if err := SendDataGoroutine(SyncRoutineBool); err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})

	//创建心跳发送协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		util.Loglevel(util.Debug, "Leader-leader", "leader-心跳发送协程启动")
		if err := SendHeartBeatGoroutine(SyncRoutineBool); err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})

	//监听信息
	for {
		select {
		case m := <-syncBean.UdpReceiveMessage:
			if m.Term >= syncBean.Term {
				b, err := StatusOperatorLeader(&m, SyncRoutineBool)
				if err != nil {
					return err
				}
				if b {
					return nil
				}
			}
		case <-syncBean.CloseChan:
			close(SyncRoutineBool)
			util.Loglevel(util.Debug, "Leader-sync", "Leader关闭")
			return nil
		}
	}
}

// SendDataGoroutineSnapshot
//
//	@Description: leader snapshot数据发送协程
//	@param SyncRoutineBool	同步协程
//	@return E	异常
func SendDataGoroutineSnapshot(SyncRoutineBool chan bool) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SendDataGoroutineSnapshot-leader", util.Strval(r))
		}
	}()
	for {
		select {
		case <-time.After(time.Second * time.Duration(config.Cluster.LeaderSnapshotSyncTime)):
			instances, deleteInstances, leaderInstances, e := CreateSendSyncDataSnapshot()
			if e != nil {
				return e
			}

			//迭代发送数据
			m := syncBean.NewMessageModel()
			syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
				m.SetMessageModeForSyncSnapShot(syncBean.Term, syncBean.Status,
					instances, deleteInstances, leaderInstances,
					value.IP, value.Port)
				syncBean.UdpSendMessage <- *m
				if config.Cluster.TrackEnable {
					util.Loglevel(util.Debug, "SendDataGoroutineSnapshot-leader",
						"leader数据同步-snapshot发送数据"+util.Strval(m.UDPTargetAddress))
				}
			})
		case <-SyncRoutineBool:

			//卸任关闭
			util.Loglevel(util.Info, "SendDataGoroutineSnapshot-leader",
				"leader-snapshot数据发送协程关闭")
			return nil
		case <-syncBean.CloseChan:

			//服务关闭
			util.Loglevel(util.Info, "SendDataGoroutineSnapshot-leader",
				"leader-snapshot数据发送协程关闭")
			return nil
		}
	}
}

// SendDataGoroutine
//
//	@Description:	leader 数据发送协程
//	@param SyncRoutineBool	同步协程
//	@return E	异常
func SendDataGoroutine(SyncRoutineBool chan bool) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SendDataGoroutine-leader", util.Strval(r))
		}
	}()
	for {

		//注册数据
		{
			if !syncBean.SectionMessage.RegisterChan.IsEmpty() {
				model := syncBean.SectionMessage.RegisterChan.Dequeue()
				m := syncBean.NewMessageModel()
				syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
					m.SetMessageModeForSyncAppend(syncBean.Term, syncBean.Status,
						*model, value.IP, value.Port)
					syncBean.UdpSendMessage <- *m
					if config.Cluster.TrackEnable {
						util.Loglevel(util.Debug, "SendDataGoroutineSnapshot-leader",
							"leader数据同步-RegisterChan-发送数据"+util.Strval(m.UDPTargetAddress))
					}
				})
			}
		}

		//删除数据
		{
			if !syncBean.SectionMessage.DeleteChan.IsEmpty() {
				model := syncBean.SectionMessage.DeleteChan.Dequeue()
				m := syncBean.NewMessageModel()
				syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
					m.SetMessageModeForSyncAppend(syncBean.Term, syncBean.Status,
						*model, value.IP, value.Port)
					syncBean.UdpSendMessage <- *m
					if config.Cluster.TrackEnable {
						util.Loglevel(util.Debug, "SendDataGoroutineSnapshot-leader",
							"leader数据同步-DeleteChan-发送数据"+util.Strval(m.UDPTargetAddress))
					}
				})
			}
		}

		//删除黑名单数据
		{
			model := syncBean.SectionMessage.CancelDeleteChan.Dequeue()
			m := syncBean.NewMessageModel()
			syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
				m.SetMessageModeForSyncAppend(syncBean.Term, syncBean.Status,
					*model, value.IP, value.Port)
				syncBean.UdpSendMessage <- *m
				if config.Cluster.TrackEnable {
					util.Loglevel(util.Debug, "SendDataGoroutineSnapshot-leader",
						"leader数据同步-CancelDeleteChan-发送数据"+util.Strval(m.UDPTargetAddress))
				}
			})
		}

		//领导者数据
		{
			if !syncBean.SectionMessage.LeaderChan.IsEmpty() {
				model := syncBean.SectionMessage.LeaderChan.Dequeue()
				m := syncBean.NewMessageModel()
				syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
					m.SetMessageModeForSyncAppend(syncBean.Term, syncBean.Status,
						*model, value.IP, value.Port)
					syncBean.UdpSendMessage <- *m
					if config.Cluster.TrackEnable {
						util.Loglevel(util.Debug, "SendDataGoroutineSnapshot-leader",
							"leader数据同步-LeaderChan-发送数据"+util.Strval(m.UDPTargetAddress))
					}
				})
			}
		}

		select {
		//卸任关闭
		case <-SyncRoutineBool:
			util.Loglevel(util.Info, "SendDataGoroutine-leader", "leader数据发送协程关闭")
			return nil

			//服务关闭
		case <-syncBean.CloseChan:
			util.Loglevel(util.Info, "SendDataGoroutine-leader", "leader数据发送协程关闭")
			return nil
		default:
			time.Sleep(time.Millisecond)
		}
	}
}

// SendHeartBeatGoroutine
//
//	@Description: leader 心跳发送协程
//	@return E	异常
func SendHeartBeatGoroutine(SyncRoutineBool chan bool) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SendHeartBeatGoroutine-leader", util.Strval(r))
		}
	}()
	for {
		select {

		//心跳
		case <-time.After(time.Duration(config.Cluster.LeaderHeartbeatTime) * time.Millisecond):
			m := syncBean.NewMessageModel()
			syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
				m.SetMessageModelForHeartBeat(syncBean.Term, syncBean.Status, value.IP, value.Port)
				syncBean.UdpSendMessage <- *m
				if config.Cluster.TrackEnable {
					util.Loglevel(util.Debug, "SendHeartBeatGoroutine-leader",
						"leader心跳发送协程发送数据"+util.Strval(m.UDPTargetAddress))
				}
			})

			//卸任关闭
		case <-SyncRoutineBool:
			util.Loglevel(util.Info, "SendDataGoroutine-leader", "leader心跳发送协程关闭")
			return nil

			//服务关闭
		case <-syncBean.CloseChan:
			util.Loglevel(util.Info, "SendDataGoroutine-leader", "leader心跳发送协程关闭")
			return nil
		}
	}
}
