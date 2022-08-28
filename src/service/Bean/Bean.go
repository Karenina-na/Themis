package Bean

import (
	"Themis/src/entity"
	"Themis/src/util"
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

//服务注册
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
