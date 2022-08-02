package service

import (
	"Envoy/src/config"
	"Envoy/src/entity"
	"Envoy/src/entity/util"
	"reflect"
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

func init() {
	RoutinePool = util.CreatePool(config.CoreRoutineNum, config.MaxRoutineNum)

	InstanceList = util.NewLinkList[entity.ServerModel]()
	DeleteInstanceList = util.NewLinkList[entity.ServerModel]()

	ServerModelList = make(map[string]map[string]*util.LinkList[entity.ServerModel])
	ServerModelList["default"] = make(map[string]*util.LinkList[entity.ServerModel])

	ServerModelQueue = make(chan entity.ServerModel, config.ServerModelQueueNum)

	ServerModelBeatQueue = make(chan entity.ServerModel, config.ServerModelBeatQueue)
	RoutinePool.CreateWork(Register)
}

func Register() {
	for {
		data := <-ServerModelQueue
		namespace := data.Namespace
		name := data.Colony + "::" + data.Name
		ServerModelListRWLock.Lock()
		if ServerModelList[namespace] == nil {
			ServerModelList[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
		}
		if ServerModelList[namespace][name] == nil {
			ServerModelList[namespace][name] = util.NewLinkList[entity.ServerModel]()
		}
		ServerModelList[namespace][name].Append(data)
		InstanceList.Append(data)
		RoutinePool.CreateWork(func() {
			wg := sync.WaitGroup{}
			wg.Add(1)
			go ServerBeat(data, namespace, name, &wg)
			wg.Wait()
			util.Loglevel(util.Info, util.Strval(data), "break ServerBeat")
		})
		ServerModelListRWLock.Unlock()
	}
}

func ServerBeat(model entity.ServerModel, namespace string, name string, wg *sync.WaitGroup) {
	defer wg.Done()
	start := time.Now().Unix()
	for {
		t := time.Now().Unix() - start
		if t == config.ServerBeatTime {
			ServerModelListRWLock.Lock()
			ServerModelList[namespace][name].DeleteByValue(model)
			if ServerModelList[namespace][name].IsEmpty() {
				delete(ServerModelList[namespace], name)
			}
			InstanceList.DeleteByValue(model)
			ServerModelListRWLock.Unlock()
			return
		}
		ServerModelBeatQueueLock.Lock()
		select {
		case data := <-ServerModelBeatQueue:
			if reflect.DeepEqual(model, data) {
				start = time.Now().Unix()
			} else {
				ServerModelBeatQueue <- data
			}
		default:
		}
		ServerModelBeatQueueLock.Unlock()
		time.Sleep(time.Millisecond)
	}
}
