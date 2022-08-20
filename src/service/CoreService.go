package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/entity/util"
	"Themis/src/exception"
	"reflect"
	"sync"
	"time"
)

func Register() (E any) {
	defer func() {
		E = recover()
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
		RoutinePool.CreateWork(func() (E any) {
			defer func() {
				E = recover()
			}()
			wg := sync.WaitGroup{}
			wg.Add(1)
			go ServerBeat(data, namespace, name, &wg)
			wg.Wait()
			util.Loglevel(util.Info, "ServerBeat:", "因心跳停止而删除-"+util.Strval(data))
			return nil
		}, func(Message any) {
			panic(exception.NewServicePanic("ServerBeat", "goroutine错误"+Message.(string)))
		})
		ServerModelListRWLock.Unlock()
	}
}

func ServerBeat(model entity.ServerModel, namespace string, name string, wg *sync.WaitGroup) (E any) {
	defer func() {
		E = recover()
	}()
	defer func() {
		wg.Done()
		DeleteMapper(&model)
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
