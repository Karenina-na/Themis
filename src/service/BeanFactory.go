package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/entity/util"
	"Themis/src/exception"
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

func ServerInitFactory() (E any) {
	defer func() {
		E = recover()
	}()
	RoutinePool = util.CreatePool(config.CoreRoutineNum, config.MaxRoutineNum)

	InstanceList = util.NewLinkList[entity.ServerModel]()
	DeleteInstanceList = util.NewLinkList[entity.ServerModel]()

	ServerModelList = make(map[string]map[string]*util.LinkList[entity.ServerModel])
	ServerModelList["default"] = make(map[string]*util.LinkList[entity.ServerModel])

	ServerModelQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	ServerModelBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)

	RoutinePool.CreateWork(Register, func(message any) {
		exception.HandleException(exception.NewServicePanic("Register", "goroutine错误"+util.Strval(message)))
	})
	if config.DatabaseEnable {
		if _, err := os.Stat("./db/Themis.db"); err == nil {
			if err := LoadDatabase(); err != nil {
				return err
			}
		}
		RoutinePool.CreateWork(Persistence, func(message any) {
			exception.HandleException(exception.NewServicePanic("Persistence", "goroutine错误"+util.Strval(message)))
		})
	}
	return nil
}
