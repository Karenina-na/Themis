package leader

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
)

// CreateSendSyncDataSnapshot
// @Description: 创建发送同步数据
// @return       instances  实例列表
// @return       list       删除实例列表
// @return       leaderList leader实例列表
// @return       E          错误
func CreateSendSyncDataSnapshot() (instances []entity.ServerModel,
	list []entity.ServerModel, leaderList []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSendSyncDataSnapshot-leader", util.Strval(r))
		}
	}()

	//创建领导者数据
	LeaderList := make([]entity.ServerModel, 0)
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	for _, Namespace := range Bean.Leaders.LeaderModelsList {
		for _, Leader := range Namespace {
			LeaderList = append(LeaderList, Leader)
		}
	}
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()

	//创建实例数据
	Instances := make([]entity.ServerModel, 0)
	Bean.InstanceList.Iterator(func(index int, value entity.ServerModel) {
		Instances = append(Instances, value)
	})

	//创建删除实例数据
	DeleteInstancesList := make([]entity.ServerModel, 0)
	Bean.DeleteInstanceList.Iterator(func(index int, value entity.ServerModel) {
		DeleteInstancesList = append(DeleteInstancesList, value)
	})
	return Instances, DeleteInstancesList, LeaderList, nil
}
