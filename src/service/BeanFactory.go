package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"fmt"
	"os"
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

	Bean.ServerModelList = make(map[string]map[string]*util.LinkList[entity.ServerModel])
	Bean.ServerModelList["default"] = make(map[string]*util.LinkList[entity.ServerModel])
	Bean.ServerModelListRWLock = sync.RWMutex{}

	Bean.ServerModelQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	Bean.ServerModelBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)

	Bean.Leaders = make(map[string]map[string]entity.ServerModel)
	Bean.Leaders["default"] = make(map[string]entity.ServerModel)
	Bean.LeadersRWLock = sync.RWMutex{}

	for i := 0; i < config.ServerModelHandleNum; i++ {
		Bean.RoutinePool.CreateWork(Register, func(message error) {
			exception.HandleException(message)
		})
	}
	Bean.RoutinePool.CreateWork(GetCenterStatusRoutine, func(message error) {
		exception.HandleException(message)
	})
	if config.DatabaseEnable {
		if _, err := os.Stat("./db/Themis.db"); err == nil {
			if err := LoadDatabase(); err != nil {
				return err
			}
		}
		Bean.RoutinePool.CreateWork(Persistence, func(message error) {
			exception.HandleException(message)
		})
	}
	fmt.Println(Bean.RoutinePool.CheckStatus())
	return nil
}
