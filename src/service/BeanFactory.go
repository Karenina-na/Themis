package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	Factory "Themis/src/pool"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"sync"
)

// InitServer
// @Description: 初始化服务
// @return       E 异常
func InitServer() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitServer-service", util.Strval(r))
		}
	}()
	//设置服务器状态信息异常处理
	util.SetStatusErrorHandle(func(err error) {
		exception.HandleException(exception.NewSystemError("ComputerStatusManager", util.Strval(err)))
	})

	//初始化实例存储链表
	Bean.InstanceList = util.NewLinkList[entity.ServerModel](func(a, b entity.ServerModel) bool {
		return a.IP == b.IP && a.Port == b.Port
	})
	Bean.DeleteInstanceList = util.NewLinkList[entity.ServerModel](func(a, b entity.ServerModel) bool {
		return a.IP == b.IP && a.Port == b.Port
	})

	//初始化服务器模型队列
	Bean.Servers = Bean.NewServersModel()
	Bean.Servers.ServerModelsList["default"] = make(map[string]*util.LinkList[entity.ServerModel])

	//初始化中间层队列
	Bean.ServersQueue = util.NewChanQueue[entity.ServerModel](config.ServerRegister.ServerModelQueueNum)
	Bean.ServersBeatQueue = util.NewChanQueue[entity.ServerModel](config.ServerBeat.ServerModelBeatQueue)

	//初始化领导者选举信息
	Bean.Leaders = Bean.NewLeadersModel()
	Bean.Leaders.LeaderModelsList["default"] = make(map[string]entity.ServerModel)
	Bean.Leaders.ElectionServers["default"] = make(map[string]*util.LinkList[entity.ServerModel])
	Bean.Leaders.LeaderModelsListRWLock = sync.RWMutex{}

	//创建服务注册协程
	for i := 0; i < config.ServerRegister.ServerModelHandleNum; i++ {
		Factory.RoutinePool.CreateWork(Register, func(message error) {
			exception.HandleException(message)
		})
	}

	//初始化初始服务器状态存储与协程
	Bean.CenterStatus = Bean.NewCenterStatusModel()
	Factory.RoutinePool.CreateWork(CenterStatusRoutine, func(message error) {
		exception.HandleException(message)
	})

	//初始化数据持久化协程
	if config.Persistence.PersistenceEnable {
		if err := LoadDatabase(); err != nil {
			return err
		}
		Factory.RoutinePool.CreateWork(Persistence, func(message error) {
			exception.HandleException(message)
		})
	}

	//初始化关闭检测通道
	Bean.ServiceCloseChan = make(chan struct{}, 0)
	return nil
}

// Close
// @Description: 关闭服务
// @return       E 异常
func Close() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Close-service", util.Strval(r))
		}
	}()
	//关闭检测通道
	close(Bean.ServiceCloseChan)
	//关闭服务注册传输队列
	Bean.ServersQueue.Destroy()
	//关闭服务心跳传输队列
	Bean.ServersBeatQueue.Destroy()
	return nil
}
