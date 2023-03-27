package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"strings"
)

// BlackInstance
// @Description: 拉黑实例
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
func BlackInstance(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("BlackInstance-service", util.Strval(r))
		}
	}()
	//将服务从服务列表中删除，加入黑名单
	Bean.Servers.ServerModelsListRWLock.Lock()
	Bean.DeleteInstanceList.Append(*model)
	Bean.InstanceList.DeleteByValue(*model)
	Bean.Servers.ServerModelsList[model.Namespace][model.Colony+"::"+model.Name].DeleteByValue(*model)

	//如果服务列表为空了，则删除该服务列表
	if Bean.Servers.ServerModelsList[model.Namespace][model.Colony+"::"+model.Name].IsEmpty() {
		delete(Bean.Servers.ServerModelsList[model.Namespace], model.Colony+"::"+model.Name)
	}
	if len(Bean.Servers.ServerModelsList[model.Namespace]) == 0 && model.Namespace != "default" {
		delete(Bean.Servers.ServerModelsList, model.Namespace)
		Bean.Leaders.LeaderModelsListRWLock.Lock()
		delete(Bean.Leaders.LeaderModelsList, model.Namespace)
		Bean.Leaders.LeaderModelsListRWLock.Unlock()
	}

	//如果删除的服务是领导者，则重新选举
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	leader := Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony]
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
	if model.Equal(&leader) {
		B, E := Election(model)
		if E != nil || !B {
			Bean.Servers.ServerModelsListRWLock.Unlock()
			return B, E
		}
	}
	Bean.Servers.ServerModelsListRWLock.Unlock()

	//集群同步
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.DeleteChan.Enqueue(model)
	}
	util.Loglevel(util.Debug, "BlackInstance", "删除服务-"+util.Strval(*model))
	return true, nil
}

// BlackColony
// @Description: 拉黑集群
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
func BlackColony(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("BlackColony-service", util.Strval(r))
		}
	}()
	//如果服务模型为空，返回错误
	Bean.Servers.ServerModelsListRWLock.Lock()
	if Bean.Servers.ServerModelsList[model.Namespace] == nil {
		Bean.Servers.ServerModelsListRWLock.Unlock()
		return false, nil
	}

	//迭代服务模型
	list := make([]string, 0, 100)
	for name, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		if colony == model.Colony {
			//将服务从服务列表中删除，加入黑名单
			L.Iterator(func(index int, server entity.ServerModel) {
				Bean.DeleteInstanceList.Append(server)
				Bean.InstanceList.DeleteByValue(server)
				Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].DeleteByValue(server)
				if Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].IsEmpty() {
					delete(Bean.Servers.ServerModelsList[server.Namespace], server.Colony+"::"+server.Name)
				}

				//集群同步
				if config.Cluster.ClusterEnable {
					syncBean.SectionMessage.DeleteChan.Enqueue(&server)
				}
			})
			list = append(list, name)
		}
	}
	//删除集群
	for _, name := range list {
		delete(Bean.Servers.ServerModelsList[model.Namespace], name)
	}
	//删除命名空间
	if len(Bean.Servers.ServerModelsList[model.Namespace]) == 0 && model.Namespace != "default" {
		delete(Bean.Servers.ServerModelsList, model.Namespace)
	}
	Bean.Servers.ServerModelsListRWLock.Unlock()

	//删除领导者
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	delete(Bean.Leaders.LeaderModelsList[model.Namespace], model.Colony)
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "BlackColony", "批量删除服务-"+util.Strval(model.Colony))
	return true, nil
}

// BlackNamespace
// @Description: 拉黑命名空间
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
func BlackNamespace(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("BlackNamespace-service", util.Strval(r))
		}
	}()
	//如果服务模型为空，返回错误
	Bean.Servers.ServerModelsListRWLock.Lock()
	if Bean.Servers.ServerModelsList[model.Namespace] == nil {
		Bean.Servers.ServerModelsListRWLock.Unlock()
		return false, nil
	}

	//迭代服务模型
	list := make([]string, 0, 100)
	for name, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		//将服务从服务列表中删除，加入黑名单
		L.Iterator(func(index int, server entity.ServerModel) {
			Bean.DeleteInstanceList.Append(server)
			Bean.InstanceList.DeleteByValue(server)
			Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].DeleteByValue(server)
			if Bean.Servers.ServerModelsList[server.Namespace][server.Colony+"::"+server.Name].IsEmpty() {
				delete(Bean.Servers.ServerModelsList[server.Namespace], server.Colony+"::"+server.Name)
			}

			//集群同步
			if config.Cluster.ClusterEnable {
				syncBean.SectionMessage.DeleteChan.Enqueue(&server)
			}
		})
		list = append(list, name)
	}
	//删除集群
	for _, name := range list {
		delete(Bean.Servers.ServerModelsList[model.Namespace], name)
	}
	//删除命名空间
	delete(Bean.Servers.ServerModelsList, model.Namespace)
	Bean.Servers.ServerModelsListRWLock.Unlock()

	//删除领导者
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	delete(Bean.Leaders.LeaderModelsList, model.Namespace)
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	util.Loglevel(util.Debug, "BlackColony", "批量删除服务-"+util.Strval(model.Colony))
	return true, nil
}

// GetBlacklistServer
// @Description: 获取黑名单服务
// @return       m 黑名单服务
// @return       E 错误
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

// DeleteInstanceFromBlacklist
// @Description: 从黑名单中删除服务
// @param        model 服务模型
// @return       B     是否成功
// @return       E     错误
func DeleteInstanceFromBlacklist(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteInstanceFromBlacklist-service", util.Strval(r))
		}
	}()
	Bean.DeleteInstanceList.DeleteByValue(*model)
	//集群同步
	if config.Cluster.ClusterEnable {
		syncBean.SectionMessage.CancelDeleteChan.Enqueue(model)
	}
	util.Loglevel(util.Debug, "DeleteInstanceFromBlacklist", "从黑名单恢复-"+util.Strval(*model))
	return true, nil
}

// GetNamespaces
//
//	@Description: 获取命名空间
//	@return n	命名空间
//	@return E	错误
func GetNamespaces() (n []string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetNamespaces-service", util.Strval(r))
		}
	}()
	list := make([]string, 0, 100)
	Bean.Servers.ServerModelsListRWLock.RLock()
	for name := range Bean.Servers.ServerModelsList {
		list = append(list, name)
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return list, nil
}

// GetColonyByNamespace
//
//	@Description: 获取集群
//	@param namespace 命名空间
//	@return c	集群
//	@return E	错误
func GetColonyByNamespace(namespace string) (c []string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetColonyByNamespace-service", util.Strval(r))
		}
	}()
	list := make([]string, 0, 100)
	Bean.Servers.ServerModelsListRWLock.RLock()

	//判断命名空间是否存在
	if Bean.Servers.ServerModelsList[namespace] == nil {
		Bean.Servers.ServerModelsListRWLock.RUnlock()
		return list, nil
	}

	//获取集群
	for name := range Bean.Servers.ServerModelsList[namespace] {
		colonyName := strings.Split(name, "::")[0]
		list = append(list, colonyName)
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return list, nil
}

// GetColonyAndServerByNamespace
//
//	@Description: 获取集群和服务
//	@param namespace 命名空间
//	@return c	集群
//	@return E	错误
func GetColonyAndServerByNamespace(namespace string) (c map[string][]string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetColonyAndServerByNamespace-service", util.Strval(r))
		}
	}()
	list := make(map[string][]string)
	Bean.Servers.ServerModelsListRWLock.RLock()

	//判断命名空间是否存在
	if Bean.Servers.ServerModelsList[namespace] == nil {
		Bean.Servers.ServerModelsListRWLock.RUnlock()
		return list, nil
	}

	//获取集群和服务名
	for name := range Bean.Servers.ServerModelsList[namespace] {
		colonyName := strings.Split(name, "::")[0]
		serverName := strings.Split(name, "::")[1]
		if list[colonyName] == nil {
			list[colonyName] = make([]string, 0)
		}
		list[colonyName] = append(list[colonyName], serverName)
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return list, nil
}

// GetInstances
// @Description: 获取所有服务实例
// @return       m 服务实例
// @return       E 错误
func GetInstances() (m map[string]map[string]map[string][]entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstances-service", util.Strval(r))
		}
	}()
	ServerLists := make(map[string]map[string]map[string][]entity.ServerModel)
	Bean.Servers.ServerModelsListRWLock.RLock()

	//迭代命名空间
	for namespace, colonyMap := range Bean.Servers.ServerModelsList {
		//如果命名空间不存在则创建
		if ServerLists[namespace] == nil {
			ServerLists[namespace] = make(map[string]map[string][]entity.ServerModel)
		}
		//迭代集群
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
			//迭代服务
			L.Iterator(func(index int, server entity.ServerModel) {
				ServerLists[namespace][colony][serverName] = append(ServerLists[namespace][colony][serverName], server)
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return ServerLists, nil
}

// GetInstance
//
//	@Description: 获取指定服务
//	@param model	服务模型
//	@return m	服务实例
//	@return E	错误
func GetInstance(model *entity.ServerModel) (L []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstance-service", util.Strval(r))
		}
	}()
	//判断是否为空
	if model.Namespace == "" || model.Colony == "" || model.Name == "" {
		return nil, exception.NewUserError("GetInstance-service", "请求非法")
	}

	ServerLists := make([]entity.ServerModel, 0)
	Bean.Servers.ServerModelsListRWLock.RLock()

	//迭代集群
	for name, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		serverName := str[1]
		if colony == model.Colony && serverName == model.Name {
			//迭代服务
			L.Iterator(func(index int, server entity.ServerModel) {
				ServerLists = append(ServerLists, server)
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return ServerLists, nil
}

// GetInstancesByNamespaceAndColony
// @Description: 获取指定命名空间和集群的服务实例--返回列表
// @param        model 服务模型
// @return       m     服务实例
// @return       E     错误
func GetInstancesByNamespaceAndColony(model *entity.ServerModel) (m []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstancesByNamespaceAndColony-service", util.Strval(r))
		}
	}()
	Bean.Servers.ServerModelsListRWLock.RLock()

	//如果命名空间为通配符则获取所有服务
	if model.Namespace == "" {
		var list []entity.ServerModel
		for _, colonyMap := range Bean.Servers.ServerModelsList {
			for _, L := range colonyMap {
				//迭代服务
				L.Iterator(func(index int, server entity.ServerModel) {
					list = append(list, server)
				})
			}
		}
		Bean.Servers.ServerModelsListRWLock.RUnlock()
		return list, nil
	}

	//如果集群空间为通配符，则获取指定命名空间下的所有服务
	if model.Colony == "" {
		var list []entity.ServerModel
		for namespace, colonyMap := range Bean.Servers.ServerModelsList {
			if namespace == model.Namespace {
				for _, L := range colonyMap {
					//迭代服务
					L.Iterator(func(index int, server entity.ServerModel) {
						list = append(list, server)
					})
				}
			}
		}
		Bean.Servers.ServerModelsListRWLock.RUnlock()
		return list, nil
	}

	//获取指定命名空间和集群的服务
	var list []entity.ServerModel
	for namespace, colonyMap := range Bean.Servers.ServerModelsList {
		if namespace == model.Namespace {
			for name, L := range colonyMap {
				str := strings.Split(name, "::")
				colony := str[0]
				if colony == model.Colony {
					//迭代服务
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

// GetInstanceByNamespaceAndColony
//
//	@Description: 获取指定命名空间和集群的服务实例--返回map
//	@param model	服务模型
//	@return m	服务实例
//	@return E	错误
func GetInstanceByNamespaceAndColony(model *entity.ServerModel) (m map[string][]entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetInstanceByNamespaceAndColony-service", util.Strval(r))
		}
	}()
	ServerLists := make(map[string][]entity.ServerModel)
	Bean.Servers.ServerModelsListRWLock.RLock()

	//迭代集群
	for name, L := range Bean.Servers.ServerModelsList[model.Namespace] {
		str := strings.Split(name, "::")
		colony := str[0]
		serverName := str[1]
		if colony == model.Colony {
			//迭代服务
			L.Iterator(func(index int, server entity.ServerModel) {
				ServerLists[serverName] = append(ServerLists[serverName], server)
			})
		}
	}
	Bean.Servers.ServerModelsListRWLock.RUnlock()
	return ServerLists, nil
}

// GetCenterStatus
// @Description: 获取中心状态
// @return       C 中心消息
// @return       E 错误
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

// GetClusterLeader
// @Description: 获取集群leader
// @return       name 集群leader
// @return       E    错误
func GetClusterLeader() (name string, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetClusterLeader-service", util.Strval(r))
		}
	}()
	if syncBean.Status != syncBean.CANDIDATE {
		return syncBean.Leader.LeaderName, nil
	}
	return "", nil
}

// GetClusterStatus
// @Description: 获取集群状态
// @return       s 集群状态
// @return       E 错误
func GetClusterStatus() (s syncBean.StatusLevel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("GetClusterStatus-service", util.Strval(r))
		}
	}()
	return syncBean.Status, nil
}
