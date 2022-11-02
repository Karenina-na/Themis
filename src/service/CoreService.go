package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"time"
)

// Register
// @Description: 注册服务
// @return       E error
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
		case <-Bean.ServiceCloseChan:
			util.Loglevel(util.Debug, "Register", "注册协程退出")
			return nil
		default:
			var data *entity.ServerModel
			flag := false
			Bean.ServersQueue.Operate(func() {
				if !Bean.ServersQueue.IsEmpty() {
					data = Bean.ServersQueue.Dequeue()
					flag = true
				}
			})
			if !flag {
				time.After(time.Millisecond)
				continue
			}
			namespace := data.Namespace
			name := data.Colony + "::" + data.Name
			Bean.Servers.ServerModelsListRWLock.Lock()
			if Bean.Servers.ServerModelsList[namespace] == nil {
				Bean.Servers.ServerModelsList[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
			}
			if Bean.Servers.ServerModelsList[namespace][name] == nil {
				Bean.Servers.ServerModelsList[namespace][name] = util.NewLinkList[entity.ServerModel]()
			}
			Bean.Servers.ServerModelsList[namespace][name].Append(*data)
			Bean.Servers.ServerModelsListRWLock.Unlock()
			Bean.Leaders.LeaderModelsListRWLock.Lock()
			if Bean.Leaders.LeaderModelsList[namespace] == nil {
				Bean.Leaders.LeaderModelsList[namespace] = make(map[string]entity.ServerModel)
			}
			if Bean.Leaders.ElectionServers[namespace] == nil {
				Bean.Leaders.ElectionServers[namespace] = make(map[string]*util.LinkList[entity.ServerModel])
			}
			if Bean.Leaders.ElectionServers[namespace][data.Colony] == nil {
				Bean.Leaders.ElectionServers[namespace][data.Colony] = util.NewLinkList[entity.ServerModel]()
			}
			Bean.Leaders.LeaderModelsListRWLock.Unlock()
			if config.ServerBeat.ServerModelBeatEnable {
				Bean.RoutinePool.CreateWork(func() (E error) {
					defer func() {
						r := recover()
						if r != nil {
							E = exception.NewUserError("BeatServer-goroutine-service", util.Strval(r))
						}
					}()
					E = ServerBeat(*data, namespace, name)
					if E != nil {
						return E
					}
					return nil
				}, func(Message error) {
					exception.HandleException(Message)
				})
			}
		}
	}
}

// ServerBeat
// @Description: 服务心跳
// @param        model     服务模型
// @param        namespace 命名空间
// @param        name      服务名称
// @return       E         error
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
		case <-Bean.ServiceCloseChan:
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
			Bean.ServersBeatQueue.Operate(func() {
				data := Bean.ServersBeatQueue.Head()
				if model.Equal(&data) {
					start = time.Now().Unix()
				}
				Bean.ServersBeatQueue.Dequeue()
			})
		}
	}
}

// CenterStatusRoutine
// @Description: 中心状态协程
// @return       E error
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
		case <-Bean.ServiceCloseChan:
			util.Loglevel(util.Debug, "CenterStatusRoutine", "监控中心协程退出")
			return nil
		case <-time.After(time.Second * time.Duration(config.ListenTime)):
			Bean.CenterStatus.CenterStatusInfoLock.Lock()
			activeNum, jobNum := Bean.RoutinePool.CheckStatus()
			Bean.CenterStatus.CenterStatusInfo.SetComputerInfoModel(util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(),
				util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
			Bean.CenterStatus.CenterStatusInfoLock.Unlock()
		}
	}
}
