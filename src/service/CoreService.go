package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"reflect"
	"time"
)

// Register 注册服务
func Register() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Register-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Register", "创建注册协程")

	for {
		data := <-ServerModelQueue
		namespace := data.Namespace
		name := data.Colony + "::" + data.Name
		LeadersRWLock.Lock()
		if Leaders[namespace] == nil {
			Leaders[namespace] = make(map[string]entity.ServerModel)
		}
		LeadersRWLock.Unlock()
		ServerModelListRWLock.Lock()
		if ServerModelList[namespace] == nil {
			ServerModelList[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
		}
		if ServerModelList[namespace][name] == nil {
			ServerModelList[namespace][name] = util.NewLinkList[entity.ServerModel]()
		}
		ServerModelList[namespace][name].Append(data)
		InstanceList.Append(data)
		RoutinePool.CreateWork(func() (E error) {
			defer func() {
				r := recover()
				if r != nil {
					E = exception.NewUserError("BeatServer-goroutine-service", util.Strval(r))
				}
			}()
			E = ServerBeat(data, namespace, name)
			if E != nil {
				return E
			}
			util.Loglevel(util.Info, "ServerBeat", "因心跳停止而删除-"+util.Strval(data))
			return nil
		}, func(Message error) {
			exception.HandleException(Message)
		})
		ServerModelListRWLock.Unlock()
	}
}

// ServerBeat 心跳服务
func ServerBeat(model entity.ServerModel, namespace string, name string) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("ServerBeat-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "ServerBeat", "创建心跳协程-"+util.Strval(model))
	start := time.Now().Unix()
	for {
		t := time.Now().Unix() - start
		if t == config.ServerBeatTime {
			ServerModelListRWLock.Lock()
			ServerModelList[namespace][name].DeleteByValue(model)
			if ServerModelList[namespace][name].IsEmpty() {
				delete(ServerModelList[namespace], name)
			}
			if len(ServerModelList[namespace]) == 0 && namespace != "default" {
				delete(ServerModelList, namespace)
				LeadersRWLock.Lock()
				delete(Leaders, namespace)
				LeadersRWLock.Unlock()
			}
			InstanceList.DeleteByValue(model)
			ServerModelListRWLock.Unlock()
			return nil
		}
		select {
		case data := <-ServerModelBeatQueue:
			if reflect.DeepEqual(model, data) {
				start = time.Now().Unix()
			} else {
				ServerModelBeatQueue <- data
			}
		default:
		}
		time.Sleep(time.Millisecond)
	}
}
