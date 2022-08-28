package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"reflect"
	"strings"
)

// DeleteServer 删除服务
func DeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-service", util.Strval(r))
		}
	}()
	Bean.ServerModelListRWLock.Lock()
	Bean.DeleteInstanceList.Append(*model)
	Bean.InstanceList.DeleteByValue(*model)
	Bean.ServerModelList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if Bean.ServerModelList[model.Namespace][model.Colony+"::"+model.Name].IsEmpty() {
		delete(Bean.ServerModelList[model.Namespace], model.Colony+"::"+model.Name)
	}
	if len(Bean.ServerModelList[model.Namespace]) == 0 && model.Name != "default" {
		delete(Bean.ServerModelList, model.Namespace)
		Bean.LeadersRWLock.Lock()
		delete(Bean.Leaders, model.Namespace)
		Bean.LeadersRWLock.Unlock()
	}
	Bean.LeadersRWLock.RLock()
	if reflect.DeepEqual(*model, Bean.Leaders[model.Namespace][model.Colony]) {
		Bean.LeadersRWLock.RUnlock()
		Bean.ServerModelListRWLock.Unlock()
		_, E := Election(model)
		if E != nil {
			return false, E
		}
	} else {
		Bean.LeadersRWLock.RUnlock()
		Bean.ServerModelListRWLock.Unlock()
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
	Bean.ServerModelListRWLock.Lock()
	if Bean.ServerModelList[model.Namespace] == nil {
		Bean.ServerModelListRWLock.Unlock()
		return false, nil
	}
	for name, L := range Bean.ServerModelList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			L.Iterator(func(index int, server entity.ServerModel) {
				Bean.DeleteInstanceList.Append(server)
				Bean.InstanceList.DeleteByValue(server)
				Bean.ServerModelList[server.Namespace][server.Colony+"::"+server.Name].DeleteByValue(server)
				if Bean.ServerModelList[server.Namespace][server.Colony+"::"+server.Name].IsEmpty() {
					delete(Bean.ServerModelList[server.Namespace], server.Colony+"::"+server.Name)
				}
			})
			list = append(list, name)
		}
	}
	for _, name := range list {
		delete(Bean.ServerModelList[model.Namespace], name)
	}
	if len(Bean.ServerModelList[model.Namespace]) == 0 && model.Name != "default" {
		delete(Bean.ServerModelList, model.Namespace)
	}
	Bean.ServerModelListRWLock.Unlock()
	Bean.LeadersRWLock.Lock()
	delete(Bean.Leaders, model.Namespace)
	Bean.LeadersRWLock.Unlock()
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
	for i := 0; i < Bean.DeleteInstanceList.Length(); i++ {
		list = append(list, Bean.DeleteInstanceList.Get(i))
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
	Bean.DeleteInstanceList.DeleteByValue(*model)
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
	Bean.ServerModelListRWLock.RLock()
	for namespace, colonyMap := range Bean.ServerModelList {
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
			L.Iterator(func(index int, server entity.ServerModel) {
				ServerLists[namespace][colony][serverName] = append(ServerLists[namespace][colony][serverName], server)
			})
		}
	}
	Bean.ServerModelListRWLock.RUnlock()
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
	Bean.ServerModelListRWLock.RLock()
	if model.Namespace == "" {
		var list []entity.ServerModel
		for _, colonyMap := range Bean.ServerModelList {
			for _, L := range colonyMap {
				L.Iterator(func(index int, server entity.ServerModel) {
					list = append(list, server)
				})
			}
		}
		return list, nil
	}
	if model.Colony == "" {
		var list []entity.ServerModel
		for namespace, colonyMap := range Bean.ServerModelList {
			if namespace == model.Namespace {
				for _, L := range colonyMap {
					L.Iterator(func(index int, server entity.ServerModel) {
						list = append(list, server)
					})
				}
			}
		}
		return list, nil
	}
	var list []entity.ServerModel
	for namespace, colonyMap := range Bean.ServerModelList {
		if namespace == model.Namespace {
			for name, L := range colonyMap {
				str := strings.Split(name, "::")
				colony := str[0]
				if colony == model.Colony {
					L.Iterator(func(index int, server entity.ServerModel) {
						list = append(list, server)
					})
				}
			}
		}
	}
	Bean.ServerModelListRWLock.RUnlock()
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
	Bean.CenterStatusLock.RLock()
	defer Bean.CenterStatusLock.RUnlock()
	return Bean.CenterStatus, nil
}
