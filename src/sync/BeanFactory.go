package sync

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
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
	syncBean.Name = config.Cluster.ClusterName
	syncBean.CloseChan = make(chan bool)
	syncBean.Status = syncBean.FOLLOW
	syncBean.Term = 1
	syncBean.Leader = *syncBean.NewLeaderModel()
	syncBean.SyncAddress = util.NewLinkList[syncBean.SyncAddressModel]()
	for _, value := range config.Cluster.Clusters {
		syncBean.SyncAddress.Append(syncBean.SyncAddressModel{IP: value["ip"], Port: value["port"]})
	}
	syncBean.UdpReceiveMessage = make(chan syncBean.MessageModel, config.Cluster.UDPQueueNum)
	syncBean.UdpSendMessage = make(chan syncBean.MessageModel, config.Cluster.UDPQueueNum)

	Bean.RoutinePool.CreateWork(func() (E error) {
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
	Bean.RoutinePool.CreateWork(func() (E error) {
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

	syncBean.SectionMessage = struct {
		RegisterChan     chan entity.ServerModel
		DeleteChan       chan entity.ServerModel
		CancelDeleteChan chan entity.ServerModel
		LeaderChan       chan entity.ServerModel
	}{
		RegisterChan:     make(chan entity.ServerModel, 10),
		DeleteChan:       make(chan entity.ServerModel, 10),
		CancelDeleteChan: make(chan entity.ServerModel, 10),
		LeaderChan:       make(chan entity.ServerModel, 10),
	}

	Bean.RoutinePool.CreateWork(func() (E error) {
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
	close(syncBean.CloseChan)
	return nil
}
