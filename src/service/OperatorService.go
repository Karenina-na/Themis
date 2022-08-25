package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"reflect"
	"strings"
)

func CheckServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckServer-service", util.Strval(r))
		}
	}()
	return InstanceList.Contain(*model), nil
}

func CheckDeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	return DeleteInstanceList.Contain(*model), nil
}

func CheckLeader(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckLeader-service", util.Strval(r))
		}
	}()
	return reflect.DeepEqual(*model, Leaders[model.Namespace]), nil
}

func DeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.Lock()
	defer ServerModelListRWLock.Unlock()
	DeleteInstanceList.Append(*model)
	InstanceList.DeleteByValue(*model)
	ServerModelList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if ServerModelList[model.Namespace][model.Colony+"::"+model.Name].IsEmpty() {
		delete(ServerModelList[model.Namespace], model.Colony+"::"+model.Name)
	}
	if ServerModelList[model.Namespace] == nil && model.Name != "default" {
		delete(ServerModelList, model.Namespace)
	}
	if reflect.DeepEqual(*model, Leaders[model.Namespace]) {
		_, E := Election(model)
		if E != nil {
			return false, E
		}
	}
	util.Loglevel(util.Debug, "DeleteServer", "删除服务-"+util.Strval(*model))
	return true, nil
}

func DeleteColonyServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteColonyServer-service", util.Strval(r))
		}
	}()
	flag := false
	list := make([]string, 0, 100)
	ServerModelListRWLock.Lock()
	defer ServerModelListRWLock.Unlock()
	if ServerModelList[model.Namespace] == nil {
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
				if reflect.DeepEqual(Leaders[server.Namespace], server) {
					flag = true
				}
			}
			list = append(list, name)
		}
	}
	for _, name := range list {
		delete(ServerModelList[model.Namespace], name)
	}
	if flag {
		_, err := Election(model)
		if err != nil {
			return false, err
		}
	}
	if ServerModelList[model.Namespace] == nil && model.Name != "default" {
		delete(ServerModelList, model.Namespace)
	}
	util.Loglevel(util.Debug, "DeleteColonyServer", "批量删除服务-"+util.Strval(model.Colony))
	return true, nil
}

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

func GetInstances() (m map[string]map[string]map[string][]entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstances-service", util.Strval(r))
		}
	}()
	ServerLists := make(map[string]map[string]map[string][]entity.ServerModel)
	ServerModelListRWLock.RLock()
	defer ServerModelListRWLock.RUnlock()
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
	return ServerLists, nil
}

func GetInstancesByNamespaceAndColony(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstancesByNamespaceAndColony-service", util.Strval(r))
		}
	}()
	ServerModelListRWLock.RLock()
	defer ServerModelListRWLock.RUnlock()
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
	return list, nil
}

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
