package Bean

import (
	"Themis/src/entity"
	"Themis/src/util"
)

// InstanceList 实例列表
var (
	// InstanceList 实力列表
	InstanceList *util.LinkList[entity.ServerModel]
	// DeleteInstanceList 实例黑名单列表
	DeleteInstanceList *util.LinkList[entity.ServerModel]
)

// Servers 服务模型
var Servers *ServersModel

// ServersQueue	服务注册通道
var ServersQueue chan entity.ServerModel

// ServersBeatQueue 服务心跳通道
var ServersBeatQueue chan entity.ServerModel

// Leaders 记账人
var Leaders *LeadersModel

// RoutinePool goroutine池
var RoutinePool *util.Pool

// CenterStatus 服务状态
var CenterStatus *CenterStatusModel

// CLOSE 服务关闭
var CLOSE chan struct{}
