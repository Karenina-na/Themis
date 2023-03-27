package sync

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	Factory "Themis/src/pool"
	"Themis/src/sync/common"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// InitSync
// @Description: 初始化同步
// @return       E error
func InitSync() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitSync-sync", util.Strval(r))
		}
	}()

	//初始化本地服务集群标志
	syncBean.Name = config.Cluster.ClusterName
	syncBean.CloseChan = make(chan bool)
	syncBean.Status = syncBean.FOLLOW
	syncBean.Term = 1
	syncBean.Leader = *syncBean.NewLeaderModel()
	syncBean.SyncAddress = util.NewLinkList[syncBean.SyncAddressModel](func(a, b syncBean.SyncAddressModel) bool {
		return a.IP == b.IP && a.Port == b.Port
	})

	//添加本地服务集群地址
	for _, value := range config.Cluster.Clusters {
		syncBean.SyncAddress.Append(syncBean.SyncAddressModel{IP: value["ip"], Port: value["port"]})
	}

	//初始化udp服务
	syncBean.UdpReceiveMessage = make(chan syncBean.MessageModel, config.Cluster.UDPQueueNum)
	syncBean.UdpSendMessage = make(chan syncBean.MessageModel, config.Cluster.UDPQueueNum)

	//创建udp接收服务
	Factory.RoutinePool.CreateWork(func() (E error) {
		defer func() {
			r := recover()
			if r != nil {
				E = exception.NewSystemError("UDPReceive-sync-goroutine", util.Strval(r))
			}
		}()
		util.Loglevel(util.Debug, "UDPReceive-sync-goroutine", "UDP接收协程启动")
		err := common.UDPReceive()
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})

	//创建udp发送服务
	Factory.RoutinePool.CreateWork(func() (E error) {
		defer func() {
			r := recover()
			if r != nil {
				E = exception.NewSystemError("UDPSend-sync-goroutine", util.Strval(r))
			}
		}()
		util.Loglevel(util.Debug, "UDPSend-sync-goroutine", "UDP发送协程启动")
		err := common.UDPSend()
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})

	//创建数据同步通道
	syncBean.SectionMessage = struct {
		RegisterChan     *util.ChanQueue[entity.ServerModel]
		DeleteChan       *util.ChanQueue[entity.ServerModel]
		CancelDeleteChan *util.ChanQueue[entity.ServerModel]
		LeaderChan       *util.ChanQueue[entity.ServerModel]
	}{
		RegisterChan:     util.NewChanQueue[entity.ServerModel](config.Cluster.LeaderQueueNum),
		DeleteChan:       util.NewChanQueue[entity.ServerModel](config.Cluster.LeaderQueueNum),
		CancelDeleteChan: util.NewChanQueue[entity.ServerModel](config.Cluster.LeaderQueueNum),
		LeaderChan:       util.NewChanQueue[entity.ServerModel](config.Cluster.LeaderQueueNum),
	}

	//创建状态控制器
	Factory.RoutinePool.CreateWork(func() (E error) {
		util.Loglevel(util.Debug, "StatusController-sync-goroutine", "状态控制器启动")
		if err := StatusController(); err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(Message)
	})
	return nil
}

// Close
// @Description: 关闭同步
// @return       E error
func Close() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Close-sync", util.Strval(r))
		}
	}()
	syncBean.SectionMessage.DeleteChan.Destroy()
	syncBean.SectionMessage.CancelDeleteChan.Destroy()
	syncBean.SectionMessage.LeaderChan.Destroy()
	syncBean.SectionMessage.RegisterChan.Destroy()
	close(syncBean.CloseChan)
	return nil
}
