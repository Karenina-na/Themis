package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"sync"
)

// InitServer 初始化服务
func InitServer() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitServer-service", util.Strval(r))
		}
	}()
	Bean.RoutinePool = util.CreatePool(config.CoreRoutineNum, config.MaxRoutineNum, config.RoutineTimeOut)
	Bean.RoutinePool.SetExceptionFunc(func(r any) {
		exception.HandleException(exception.NewSystemError("Pool池", util.Strval(r)))
	})

	util.SetStatusErrorHandle(func(err error) {
		exception.HandleException(exception.NewSystemError("ComputerStatusManager", util.Strval(err)))
	})

	Bean.InstanceList = util.NewLinkList[entity.ServerModel]()
	Bean.DeleteInstanceList = util.NewLinkList[entity.ServerModel]()

	Bean.Servers = Bean.NewServersModel()
	Bean.Servers.ServerModelsList["default"] = make(map[string]*util.LinkList[entity.ServerModel])

	Bean.ServersQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	Bean.ServersBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)

	Bean.Leaders = Bean.NewLeadersModel()
	Bean.Leaders.LeaderModelsList["default"] = make(map[string]entity.ServerModel)
	Bean.Leaders.LeaderModelsListRWLock = sync.RWMutex{}

	for i := 0; i < config.ServerModelHandleNum; i++ {
		Bean.RoutinePool.CreateWork(Register, func(message error) {
			exception.HandleException(message)
		})
	}
	Bean.CenterStatus = Bean.NewCenterStatusModel()
	Bean.RoutinePool.CreateWork(CenterStatusRoutine, func(message error) {
		exception.HandleException(message)
	})
	if config.DatabaseEnable {
		if err := LoadDatabase(); err != nil {
			return err
		}
		Bean.RoutinePool.CreateWork(Persistence, func(message error) {
			exception.HandleException(message)
		})
	}
	Bean.CLOSE = make(chan struct{}, 0)
	return nil
}

// Close 关闭服务
func Close() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Close-service", util.Strval(r))
		}
	}()
	close(Bean.CLOSE)
	Bean.RoutinePool.Close()
	return nil
}
