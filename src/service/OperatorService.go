package service

import (
	"Themis/src/entity"
	"reflect"
	"strings"
)

func CheckServer(model *entity.ServerModel) bool {
	return InstanceList.Contain(*model)
}

func CheckDeleteServer(model *entity.ServerModel) bool {
	return DeleteInstanceList.Contain(*model)
}

func CheckLeader(model *entity.ServerModel) bool {
	return reflect.DeepEqual(*model, Leader)
}

func DeleteServer(model *entity.ServerModel) bool {
	ServerModelListRWLock.Lock()
	defer ServerModelListRWLock.Unlock()
	DeleteInstanceList.Append(*model)
	Assert := InstanceList.DeleteByValue(*model) &&
		ServerModelList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if reflect.DeepEqual(*model, Leader) {
		Election(model)
	}
	return Assert
}

func DeleteColony(model *entity.ServerModel) bool {
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
		Election(model)
	}
	return true
}

func GetDeleteInstances() []entity.ServerModel {
	list := make([]entity.ServerModel, 0, 100)
	for i := 0; i < DeleteInstanceList.Length(); i++ {
		list = append(list, DeleteInstanceList.Get(i))
	}
	return list
}

func DeleteDeleteInstance(model *entity.ServerModel) bool {
	DeleteInstanceList.DeleteByValue(*model)
	return true
}

func GetInstances() map[string]map[string]map[string][]entity.ServerModel {
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
	return ServerLists
}
