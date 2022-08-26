package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"reflect"
	"strings"
)

// CheckServer 检查服务是否存在
func CheckServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckServer-service", util.Strval(r))
		}
	}()
	return InstanceList.Contain(*model), nil
}

// CheckDeleteServer 检查服务是否存在黑名单
func CheckDeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	return DeleteInstanceList.Contain(*model), nil
}

// CheckLeader 检查是否是领导
func CheckLeader(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckLeader-service", util.Strval(r))
		}
	}()
	LeadersRWLock.RLock()
	Assert := reflect.DeepEqual(*model, Leaders[model.Namespace][model.Colony])
	LeadersRWLock.RUnlock()
	return Assert, nil
}

// DeleteServer 删除服务
func DeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.Lock()
	DeleteInstanceList.Append(*model)
	InstanceList.DeleteByValue(*model)
	ServerModelList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if ServerModelList[model.Namespace][model.Colony+"::"+model.Name].IsEmpty() {
		delete(ServerModelList[model.Namespace], model.Colony+"::"+model.Name)
	}
	if len(ServerModelList[model.Namespace]) == 0 && model.Name != "default" {
		delete(ServerModelList, model.Namespace)
		LeadersRWLock.Lock()
		delete(Leaders, model.Namespace)
		LeadersRWLock.Unlock()
	}
	LeadersRWLock.RLock()
	if reflect.DeepEqual(*model, Leaders[model.Namespace][model.Colony]) {
		LeadersRWLock.RUnlock()
		ServerModelListRWLock.Unlock()
		_, E := Election(model)
		if E != nil {
			return false, E
		}
	} else {
		LeadersRWLock.RUnlock()
		ServerModelListRWLock.Unlock()
	}
	util.Loglevel(util.Debug, "DeleteServer", "删除服务-"+util.Strval(*model))
	return true, nil
}

// DeleteColonyServer 删除集群服务
func DeleteColonyServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteColonyServer-service", util.Strval(r))
		}
	}()
	list := make([]string, 0, 100)
	ServerModelListRWLock.Lock()
	if ServerModelList[model.Namespace] == nil {
		ServerModelListRWLock.Unlock()
		return false, nil
	}
	for name, L := range ServerModelList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			for i := 0; i < L.Length(); i++ {
				server := L.Get(i)
				DeleteInstanceList.Append(server)
				InstanceList.DeleteByValue(server)
				ServerModelList[server.Namespace][server.Colony+"::"+server.Name].DeleteByValue(server)
				if ServerModelList[server.Namespace][server.Colony+"::"+server.Name].IsEmpty() {
					delete(ServerModelList[server.Namespace], server.Colony+"::"+server.Name)
				}
			}
			list = append(list, name)
		}
	}
	for _, name := range list {
		delete(ServerModelList[model.Namespace], name)
	}
	if len(ServerModelList[model.Namespace]) == 0 && model.Name != "default" {
		delete(ServerModelList, model.Namespace)
	}
	ServerModelListRWLock.Unlock()
	LeadersRWLock.Lock()
	delete(Leaders, model.Namespace)
	LeadersRWLock.Unlock()
	util.Loglevel(util.Debug, "DeleteColonyServer", "批量删除服务-"+util.Strval(model.Colony))
	return true, nil
}

// GetBlacklistServer 获取黑名单服务
func GetBlacklistServer() (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetBlacklistServer-service", util.Strval(r))
		}
	}()
	list := make([]entity.ServerModel, 0, 100)
	for i := 0; i < DeleteInstanceList.Length(); i++ {
		list = append(list, DeleteInstanceList.Get(i))
	}
	return list, nil
}

// DeleteInstanceFromBlacklist 删除黑名单中的服务
func DeleteInstanceFromBlacklist(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteInstanceFromBlacklist-service", util.Strval(r))
		}
	}()
	DeleteInstanceList.DeleteByValue(*model)
	util.Loglevel(util.Debug, "DeleteInstanceFromBlacklist", "从黑名单恢复-"+util.Strval(*model))
	return true, nil
}

// GetInstances 获取所有服务
func GetInstances() (m map[string]map[string]map[string][]entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstances-service", util.Strval(r))
		}
	}()
	ServerLists := make(map[string]map[string]map[string][]entity.ServerModel)
	ServerModelListRWLock.Lock()
	for namespace, colonyMap := range ServerModelList {
		if ServerLists[namespace] == nil {
			ServerLists[namespace] = make(map[string]map[string][]entity.ServerModel)
		}
		for name, L := range colonyMap {
			str := strings.Split(name, "::")
			colony := str[0]
			serverName := str[1]
			if ServerLists[namespace][colony] == nil {
				ServerLists[namespace][colony] = make(map[string][]entity.ServerModel)
			}
			if ServerLists[namespace][colony][serverName] == nil {
				ServerLists[namespace][colony][serverName] = make([]entity.ServerModel, 0, 100)
			}
			for i := 0; i < L.Length(); i++ {
				ServerLists[namespace][colony][serverName] = append(ServerLists[namespace][colony][serverName], entity.ServerModel{
					IP:        L.Get(i).IP,
					Name:      L.Get(i).Name,
					Port:      L.Get(i).Port,
					Time:      L.Get(i).Time,
					Colony:    L.Get(i).Colony,
					Namespace: L.Get(i).Namespace,
				})
			}
		}
	}
	ServerModelListRWLock.Unlock()
	return ServerLists, nil
}

// GetInstancesByNamespaceAndColony 获取指定命名空间下指定集群的所有服务
func GetInstancesByNamespaceAndColony(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstancesByNamespaceAndColony-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.RLock()
	if model.Namespace == "" {
		var list []entity.ServerModel
		for _, colonyMap := range ServerModelList {
			for _, L := range colonyMap {
				for i := 0; i < L.Length(); i++ {
					list = append(list, entity.ServerModel{
						IP:        L.Get(i).IP,
						Name:      L.Get(i).Name,
						Port:      L.Get(i).Port,
						Time:      L.Get(i).Time,
						Colony:    L.Get(i).Colony,
						Namespace: L.Get(i).Namespace,
					})
				}
			}
		}
		return list, nil
	}
	if model.Colony == "" {
		var list []entity.ServerModel
		for namespace, colonyMap := range ServerModelList {
			if namespace == model.Namespace {
				for _, L := range colonyMap {
					for i := 0; i < L.Length(); i++ {
						list = append(list, entity.ServerModel{
							IP:        L.Get(i).IP,
							Name:      L.Get(i).Name,
							Port:      L.Get(i).Port,
							Time:      L.Get(i).Time,
							Colony:    L.Get(i).Colony,
							Namespace: L.Get(i).Namespace,
						})
					}
				}
			}
		}
		return list, nil
	}
	var list []entity.ServerModel
	for namespace, colonyMap := range ServerModelList {
		if namespace == model.Namespace {
			for name, L := range colonyMap {
				str := strings.Split(name, "::")
				colony := str[0]
				if colony == model.Colony {
					for i := 0; i < L.Length(); i++ {
						list = append(list, entity.ServerModel{
							IP:        L.Get(i).IP,
							Name:      L.Get(i).Name,
							Port:      L.Get(i).Port,
							Time:      L.Get(i).Time,
							Colony:    L.Get(i).Colony,
							Namespace: L.Get(i).Namespace,
						})
					}
				}
			}
		}
	}
	ServerModelListRWLock.RUnlock()
	return list, nil
}

// GetCenterStatus 获取中心状态
func GetCenterStatus() (C *entity.ComputerInfoModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetCenterStatus-service", util.Strval(r))
		}
	}()
	activeNum, jobNum := RoutinePool.CheckStatus()
	computerStatus := entity.NewComputerInfoModel(
		util.GetCpuInfo(), *util.GetMemInfo(), *util.GetHostInfo(), util.GetDiskInfo(), util.GetNetInfo(), activeNum, jobNum)
	return computerStatus, nil
}
