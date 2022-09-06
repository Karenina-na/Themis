package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"sync"
)

//
// InitServer
// @Description: 初始化服务
// @return       E 异常
//
func InitServer() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitServer-service", util.Strval(r))
		}
	}()
	Bean.RoutinePool = util.CreatePool(config.Goroutine.CoreRoutineNum,
		config.Goroutine.MaxRoutineNum, config.Goroutine.RoutineTimeOut)
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

	Bean.ServersQueue = make(chan entity.ServerModel, config.ServerRegister.ServerModelQueueNum)

	Bean.ServersBeatQueue = make(chan entity.ServerModel, config.ServerBeat.ServerModelBeatQueue)

	Bean.Leaders = Bean.NewLeadersModel()
	Bean.Leaders.LeaderModelsList["default"] = make(map[string]entity.ServerModel)
	Bean.Leaders.LeaderModelsListRWLock = sync.RWMutex{}

	for i := 0; i < config.ServerRegister.ServerModelHandleNum; i++ {
		Bean.RoutinePool.CreateWork(Register, func(message error) {
			exception.HandleException(message)
		})
	}
	Bean.CenterStatus = Bean.NewCenterStatusModel()
	Bean.RoutinePool.CreateWork(CenterStatusRoutine, func(message error) {
		exception.HandleException(message)
	})
	if config.Persistence.PersistenceEnable {
		if err := LoadDatabase(); err != nil {
			return err
		}
		Bean.RoutinePool.CreateWork(Persistence, func(message error) {
			exception.HandleException(message)
		})
	}
	Bean.ServiceCloseChan = make(chan struct{}, 0)
	return nil
}

//
// Close
// @Description: 关闭服务
// @return       E 异常
//
func Close() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Close-service", util.Strval(r))
		}
	}()
	close(Bean.ServiceCloseChan)
	Bean.RoutinePool.Close()
	return nil
}
