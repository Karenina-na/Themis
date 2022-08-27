package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"os"
	"sync"
	"time"
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
)

//记账
var (
	// Leaders 记账人
	Leaders map[string]map[string]entity.ServerModel
	// LeadersRWLock 记账人读写锁
	LeadersRWLock sync.RWMutex
)

// RoutinePool goroutine池
var RoutinePool *util.Pool

var (
	// CenterStatus 服务器状态
	CenterStatus *entity.ComputerInfoModel
	// CenterStatusLock  服务器状态读写锁
	CenterStatusLock sync.RWMutex
)

// InitServer 初始化服务
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

	util.SetStatusErrorHandle(func(err error) {
		exception.HandleException(exception.NewSystemError("ComputerStatusManager", util.Strval(err)))
	})

	InstanceList = util.NewLinkList[entity.ServerModel]()
	DeleteInstanceList = util.NewLinkList[entity.ServerModel]()

	ServerModelList = make(map[string]map[string]*util.LinkList[entity.ServerModel])
	ServerModelList["default"] = make(map[string]*util.LinkList[entity.ServerModel])
	ServerModelListRWLock = sync.RWMutex{}

	ServerModelQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	ServerModelBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)

	Leaders = make(map[string]map[string]entity.ServerModel)
	Leaders["default"] = make(map[string]entity.ServerModel)
	LeadersRWLock = sync.RWMutex{}

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
	RoutinePool.CreateWork(GetCenterStatusRoutine, func(message error) {
		exception.HandleException(message)
	})
	return nil
}

func GetCenterStatusRoutine() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Register-service", util.Strval(r))
		}
	}()
	for {
		CenterStatusLock.Lock()
		activeNum, jobNum := RoutinePool.CheckStatus()
		computerStatus := entity.NewComputerInfoModel(
			util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(), util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
		CenterStatus = computerStatus
		CenterStatusLock.Unlock()
		time.Sleep(time.Second * time.Duration(config.ListenTime))
	}
}
