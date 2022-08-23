package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"reflect"
	"time"
)

func Register() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Register-service", util.Strval(r))
		}
	}()
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

func ServerBeat(model entity.ServerModel, namespace string, name string) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("ServerBeat-service", util.Strval(r))
		}
	}()
	defer func() {
		E = DeleteMapper(&model)
	}()
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
			return nil
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
