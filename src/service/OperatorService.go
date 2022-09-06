package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"reflect"
	"strings"
)

//
// DeleteServer
// @Description: 删除服务
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
//
func DeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-service", util.Strval(r))
		}
	}()
	Bean.Servers.ServerModelsListRWLock.Lock()
	Bean.DeleteInstanceList.Append(*model)
	Bean.InstanceList.DeleteByValue(*model)
	Bean.Servers.ServerModelsList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)
	if Bean.Servers.ServerModelsList[model.Namespace][model.Colony+"::"+model.Name].IsEmpty() {
		delete(Bean.Servers.ServerModelsList[model.Namespace], model.Colony+"::"+model.Name)
	}
	if len(Bean.Servers.ServerModelsList[model.Namespace]) == 0 && model.Name != "default" {
		delete(Bean.Servers.ServerModelsList, model.Namespace)
		Bean.Leaders.LeaderModelsListRWLock.Lock()
		delete(Bean.Leaders.LeaderModelsList, model.Namespace)
		Bean.Leaders.LeaderModelsListRWLock.Unlock()
	}
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	if reflect.DeepEqual(*model, Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony]) {
		Bean.Leaders.LeaderModelsListRWLock.RUnlock()
		Bean.Servers.ServerModelsListRWLock.Unlock()
		_, E := Election(model)
		if E != nil {
			return false, E
		}
	} else {
		Bean.Leaders.LeaderModelsListRWLock.RUnlock()
		Bean.Servers.ServerModelsListRWLock.Unlock()
	}
	util.Loglevel(util.Debug, "DeleteServer", "删除服务-"+util.Strval(*model))
	return true, nil
}

//
// DeleteColonyServer
// @Description: 删除集群
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
//
func DeleteColonyServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteColonyServer-service", util.Strval(r))
		}
	}()
	list := make([]string, 0, 100)
	Bean.Servers.ServerModelsListRWLock.Lock()
	if Bean.Servers.ServerModelsList[model.Namespace] == nil {
		Bean.Servers.ServerModelsListRWLock.Unlock()
		return false, nil
	}
	for name, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			L.Iterator(func(index int, server entity.ServerModel) {
				Bean.DeleteInstanceList.Append(server)
				Bean.InstanceList.DeleteByValue(server)
				Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].DeleteByValue(server)
				if Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].IsEmpty() {
					delete(Bean.Servers.ServerModelsList[server.Namespace], server.Colony+"::"+server.Name)
				}
			})
			list = append(list, name)
		}
	}
	for _, name := range list {
		delete(Bean.Servers.ServerModelsList[model.Namespace], name)
	}
	if len(Bean.Servers.ServerModelsList[model.Namespace]) == 0 && model.Name != "default" {
		delete(Bean.Servers.ServerModelsList, model.Namespace)
	}
	Bean.Servers.ServerModelsListRWLock.Unlock()
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	delete(Bean.Leaders.LeaderModelsList, model.Namespace)
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "DeleteColonyServer", "批量删除服务-"+util.Strval(model.Colony))
	return true, nil
}

//
// GetBlacklistServer
// @Description: 获取黑名单服务
// @return       m 黑名单服务
// @return       E 错误
//
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

//
// DeleteInstanceFromBlacklist
// @Description: 从黑名单中删除服务
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
//
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

//
// GetInstances
// @Description: 获取所有服务实例
// @return       m 服务实例
// @return       E 错误
//
func GetInstances() (m map[string]map[string]map[string][]entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstances-service", util.Strval(r))
		}
	}()
	ServerLists := make(map[string]map[string]map[string][]entity.ServerModel)
	Bean.Servers.ServerModelsListRWLock.RLock()
	for namespace, colonyMap := range Bean.Servers.ServerModelsList {
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
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return ServerLists, nil
}

//
// GetInstancesByNamespaceAndColony
// @Description: 获取指定命名空间和集群的服务实例
// @param        model 服务模型
// @return       m     服务实例
// @return       E     错误
//
func GetInstancesByNamespaceAndColony(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstancesByNamespaceAndColony-service", util.Strval(r))
		}
	}()
	Bean.Servers.ServerModelsListRWLock.RLock()
	if model.Namespace == "" {
		var list []entity.ServerModel
		for _, colonyMap := range Bean.Servers.ServerModelsList {
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
		for namespace, colonyMap := range Bean.Servers.ServerModelsList {
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
	for namespace, colonyMap := range Bean.Servers.ServerModelsList {
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
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return list, nil
}

//
// GetCenterStatus
// @Description: 获取中心状态
// @return       C 中心消息
// @return       E 错误
//
func GetCenterStatus() (C *entity.ComputerInfoModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetCenterStatus-service", util.Strval(r))
		}
	}()
	Bean.CenterStatus.CenterStatusInfoLock.RLock()
	defer Bean.CenterStatus.CenterStatusInfoLock.RUnlock()
	return Bean.CenterStatus.CenterStatusInfo, nil
}

//
// GetClusterLeader
// @Description: 获取集群leader
// @return       name 集群leader
// @return       E    错误
//
func GetClusterLeader() (name string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetClusterLeader-service", util.Strval(r))
		}
	}()
	if syncBean.Status != syncBean.CANDIDATE {
		return syncBean.LeaderName, nil
	}
	return "", nil
}

//
// GetClusterStatus
// @Description: 获取集群状态
// @return       s 集群状态
// @return       E 错误
//
func GetClusterStatus() (s syncBean.StatusLevel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetClusterStatus-service", util.Strval(r))
		}
	}()
	return syncBean.Status, nil
}
