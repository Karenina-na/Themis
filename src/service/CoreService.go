package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
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
		data := <-Bean.ServerModelQueue
		namespace := data.Namespace
		name := data.Colony + "::" + data.Name
		Bean.LeadersRWLock.Lock()
		if Bean.Leaders[namespace] == nil {
			Bean.Leaders[namespace] = make(map[string]entity.ServerModel)
		}
		Bean.LeadersRWLock.Unlock()
		Bean.ServerModelListRWLock.Lock()
		if Bean.ServerModelList[namespace] == nil {
			Bean.ServerModelList[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
		}
		if Bean.ServerModelList[namespace][name] == nil {
			Bean.ServerModelList[namespace][name] = util.NewLinkList[entity.ServerModel]()
		}
		Bean.ServerModelList[namespace][name].Append(data)
		Bean.InstanceList.Append(data)
		if config.ServerModelBeatEnable {
			Bean.RoutinePool.CreateWork(func() (E error) {
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
		}
		Bean.ServerModelListRWLock.Unlock()
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
			Bean.ServerModelListRWLock.Lock()
			Bean.ServerModelList[namespace][name].DeleteByValue(model)
			if Bean.ServerModelList[namespace][name].IsEmpty() {
				delete(Bean.ServerModelList[namespace], name)
			}
			if len(Bean.ServerModelList[namespace]) == 0 && namespace != "default" {
				delete(Bean.ServerModelList, namespace)
				Bean.LeadersRWLock.Lock()
				delete(Bean.Leaders, namespace)
				Bean.LeadersRWLock.Unlock()
			}
			Bean.InstanceList.DeleteByValue(model)
			Bean.ServerModelListRWLock.Unlock()
			return nil
		}
		select {
		case data := <-Bean.ServerModelBeatQueue:
			if reflect.DeepEqual(model, data) {
				start = time.Now().Unix()
			} else {
				Bean.ServerModelBeatQueue <- data
			}
		default:
		}
		time.Sleep(time.Millisecond)
	}
}

// GetCenterStatusRoutine 获取中心状态服务
func GetCenterStatusRoutine() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Register-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "GetCenterStatusRoutine", "创建监控中心协程")
	for {
		Bean.CenterStatusLock.Lock()
		activeNum, jobNum := Bean.RoutinePool.CheckStatus()
		computerStatus := entity.NewComputerInfoModel(
			util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(), util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
		Bean.CenterStatus = computerStatus
		Bean.CenterStatusLock.Unlock()
		time.Sleep(time.Second * time.Duration(config.ListenTime))
	}
}
