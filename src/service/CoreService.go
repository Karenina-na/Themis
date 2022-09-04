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
		select {
		case <-time.After(time.Second):
			time.Sleep(time.Millisecond)
		case <-Bean.CLOSE:
			util.Loglevel(util.Debug, "Register", "注册协程退出")
			return nil
		case data := <-Bean.ServersQueue:
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
			if config.ServerBeat.ServerModelBeatEnable {
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
					return nil
				}, func(Message error) {
					exception.HandleException(Message)
				})
			}
			Bean.Servers.ServerModelsListRWLock.Unlock()
		}
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
		select {
		case <-Bean.CLOSE:
			util.Loglevel(util.Debug, "ServerBeat", "心跳协程退出-"+util.Strval(model))
			return nil
		case <-time.After(time.Millisecond):
			t := time.Now().Unix() - start
			if t == config.ServerBeat.ServerBeatTime {
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
				util.Loglevel(util.Info, "ServerBeat", "因心跳停止而删除-"+util.Strval(model))
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
		}
	}
}

// CenterStatusRoutine 监控中心协程
func CenterStatusRoutine() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Register-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "CenterStatusRoutine", "创建监控中心协程")
	for {
		select {
		case <-Bean.CLOSE:
			util.Loglevel(util.Debug, "CenterStatusRoutine", "监控中心协程退出")
			return nil
		case <-time.After(time.Second * time.Duration(config.ListenTime)):
			Bean.CenterStatus.CenterStatusInfoLock.Lock()
			activeNum, jobNum := Bean.RoutinePool.CheckStatus()
			Bean.CenterStatus.CenterStatusInfo.SetComputerInfoModel(util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(),
				util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
			Bean.CenterStatus.CenterStatusInfoLock.Unlock()
		case <-time.After(time.Millisecond):
		}
	}
}
