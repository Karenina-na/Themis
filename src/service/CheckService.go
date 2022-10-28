package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
)

// CheckServer
// @Description: 检查服务是否存在
// @param        model 服务模型
// @return       B     是否存在
// @return       E     错误
func CheckServer(m *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckServer-service", util.Strval(r))
		}
	}()
	flag := false
	Bean.InstanceList.Iterator(func(index int, model entity.ServerModel) {
		flag = m.Equal(&model)
		if flag {
			return
		}
	})
	return flag, nil
}

// CheckDeleteServer
// @Description: 检查服务是否在黑名单
// @param        model 服务模型
// @return       B     是否存在
// @return       E     错误
func CheckDeleteServer(m *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	flag := false
	Bean.DeleteInstanceList.Iterator(func(index int, model entity.ServerModel) {
		flag = m.Equal(&model)
		if flag {
			return
		}
	})
	return flag, nil
}

// CheckLeader
// @Description: 检查服务是否为Leader
// @param        model 服务模型
// @return       B     是否存在
// @return       E     错误
func CheckLeader(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckLeader-service", util.Strval(r))
		}
	}()
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	serverModel := Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony].Clone()
	Assert := serverModel.Equal(model)
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
	return Assert, nil
}
