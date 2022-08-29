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
		data := <-Bean.ServersQueue
		namespace := data.Namespace
		name := data.Colony + "::" + data.Name
		Bean.Leaders.LeaderModelsListRWLock.Lock()
		if Bean.Leaders.LeaderModelsList[namespace] == nil {
			Bean.Leaders.LeaderModelsList[namespace] = make(map[string]entity.ServerModel)
		}
		Bean.Leaders.LeaderModelsListRWLock.Unlock()
		Bean.Servers.ServerModelsListRWLock.Lock()
		if Bean.Servers.ServerModelsList[namespace] == nil {
			Bean.Servers.ServerModelsList[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
		}
		if Bean.Servers.ServerModelsList[namespace][name] == nil {
			Bean.Servers.ServerModelsList[namespace][name] = util.NewLinkList[entity.ServerModel]()
		}
		Bean.Servers.ServerModelsList[namespace][name].Append(data)
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
		Bean.Servers.ServerModelsListRWLock.Unlock()
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
			Bean.Servers.ServerModelsListRWLock.Lock()
			Bean.Servers.ServerModelsList[namespace][name].DeleteByValue(model)
			if Bean.Servers.ServerModelsList[namespace][name].IsEmpty() {
				delete(Bean.Servers.ServerModelsList[namespace], name)
			}
			if len(Bean.Servers.ServerModelsList[namespace]) == 0 && namespace != "default" {
				delete(Bean.Servers.ServerModelsList, namespace)
				Bean.Leaders.LeaderModelsListRWLock.Lock()
				delete(Bean.Leaders.LeaderModelsList, namespace)
				Bean.Leaders.LeaderModelsListRWLock.Unlock()
			}
			Bean.InstanceList.DeleteByValue(model)
			Bean.Servers.ServerModelsListRWLock.Unlock()
			return nil
		}
		select {
		case data := <-Bean.ServersBeatQueue:
			if reflect.DeepEqual(model, data) {
				start = time.Now().Unix()
			} else {
				Bean.ServersBeatQueue <- data
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
		Bean.CenterStatus.CenterStatusInfoLock.Lock()
		activeNum, jobNum := Bean.RoutinePool.CheckStatus()
		Bean.CenterStatus.CenterStatusInfo.SetComputerInfoModel(util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(),
			util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
		Bean.CenterStatus.CenterStatusInfoLock.Unlock()
		time.Sleep(time.Second * time.Duration(config.ListenTime))
	}
}
