package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"os"
	"sync"
)

// InstanceList 实例列表
var (
	// InstanceList 实力列表
	InstanceList *util.LinkList[entity.ServerModel]
	// DeleteInstanceList 实例黑名单列表
	DeleteInstanceList *util.LinkList[entity.ServerModel]
)

//服务模型
var (
	// ServerModelList 服务模型
	ServerModelList map[string]map[string]*util.LinkList[entity.ServerModel]

	// ServerModelListRWLock 服务模型读写锁
	ServerModelListRWLock sync.RWMutex
)
var (
	// ServerModelQueue	服务注册通道
	ServerModelQueue chan entity.ServerModel
)

//服务心跳
var (
	// ServerModelBeatQueue 服务心跳通道
	ServerModelBeatQueue chan entity.ServerModel

	// ServerModelBeatQueueLock 服务心跳通道读写锁
	ServerModelBeatQueueLock sync.RWMutex
)

//记账
var (
	// Leader 记账人
	Leader entity.ServerModel
)

// RoutinePool goroutine池
var RoutinePool *util.Pool

func InitServer() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitServer-service", util.Strval(r))
		}
	}()
	RoutinePool = util.CreatePool(config.CoreRoutineNum, config.MaxRoutineNum)
	RoutinePool.SetExceptionFunc(func(r any) {
		exception.HandleException(exception.NewSystemError("Pool池", util.Strval(r)))
	})

	InstanceList = util.NewLinkList[entity.ServerModel]()
	DeleteInstanceList = util.NewLinkList[entity.ServerModel]()

	ServerModelList = make(map[string]map[string]*util.LinkList[entity.ServerModel])
	ServerModelList["default"] = make(map[string]*util.LinkList[entity.ServerModel])

	ServerModelQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	ServerModelBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)

	RoutinePool.CreateWork(Register, func(message error) {
		exception.HandleException(message)
	})
	if config.DatabaseEnable {
		if _, err := os.Stat("./db/Themis.db"); err == nil {
			if err := LoadDatabase(); err != nil {
				return err
			}
		}
		RoutinePool.CreateWork(Persistence, func(message error) {
			exception.HandleException(message)
		})
	}
	return nil
}
