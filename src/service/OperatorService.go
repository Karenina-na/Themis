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
	return reflect.DeepEqual(*model, Leader), nil
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
	Assert := InstanceList.DeleteByValue(*model) &&
		ServerModelList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if reflect.DeepEqual(*model, Leader) {
		_, E := Election(model)
		if E != nil {
			return false, E
		}
	}
	return Assert, nil
}

func DeleteColony(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteColony-service", util.Strval(r))
		}
	}()
	flag := false
	list := make([]string, 0, 100)
	ServerModelListRWLock.Lock()
	defer ServerModelListRWLock.Unlock()
	for name, L := range ServerModelList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			for i := 0; i < L.Length(); i++ {
				server := L.Get(i)
				DeleteInstanceList.Append(server)
				InstanceList.DeleteAllByValue(server)
				if reflect.DeepEqual(Leader, server) {
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
	return true, nil
}

func GetDeleteInstances() (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetDeleteInstances-service", util.Strval(r))
		}
	}()
	list := make([]entity.ServerModel, 0, 100)
	for i := 0; i < DeleteInstanceList.Length(); i++ {
		list = append(list, DeleteInstanceList.Get(i))
	}
	return list, nil
}

func DeleteDeleteInstance(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteDeleteInstance-service", util.Strval(r))
		}
	}()
	DeleteInstanceList.DeleteByValue(*model)
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
					Time:      L.Get(i).Time,
					Colony:    L.Get(i).Colony,
					Namespace: L.Get(i).Namespace,
				})
			}
		}
	}
	return ServerLists, nil
}

func GetCenterStatus() (A int, J int, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetCenterStatus-service", util.Strval(r))
		}
	}()
	activeNum, jobNum := RoutinePool.CheckStatus()
	return activeNum, jobNum, nil
}
