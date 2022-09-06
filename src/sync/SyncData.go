package sync

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service"
	"Themis/src/service/Bean"
	"Themis/src/util"
)

//
// DataSyncInstances
// @Description: 同步实例数据
// @param        list []entity.ServerModel 实例列表
// @return       E    error                    错误信息
//
func DataSyncInstances(list []entity.ServerModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DataSyncInstances-sync", util.Strval(r))
		}
	}()
	for _, v := range list {
		if !Bean.InstanceList.Contain(v) {
			B, err := service.RegisterServer(&v)
			if err != nil || !B {
				return err
			}
		}
	}
	return nil
}

//
// DataSyncDelete
// @Description: 同步删除数据
// @param        list []entity.ServerModel 实例列表
// @return       E    error                    错误信息
//
func DataSyncDelete(list []entity.ServerModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DataSyncDelete-sync", util.Strval(r))
		}
	}()
	for _, v := range list {
		if !Bean.DeleteInstanceList.Contain(v) {
			B, err := service.DeleteServer(&v)
			if err != nil || !B {
				return err
			}
		}
	}
	return nil
}

//
// DataSyncLeader
// @Description: 同步leader数据
// @param        list []entity.ServerModel 实例列表
// @return       E    error                    错误信息
//
func DataSyncLeader(list []entity.ServerModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DataSyncLeader-sync", util.Strval(r))
		}
	}()
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	for _, v := range list {
		Bean.Leaders.LeaderModelsList[v.Namespace][v.Colony] = v
	}
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	return nil
}
