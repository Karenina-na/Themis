package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"reflect"
)

// CheckServer 检查服务是否存在
func CheckServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckServer-service", util.Strval(r))
		}
	}()
	return Bean.InstanceList.Contain(*model), nil
}

// CheckDeleteServer 检查服务是否存在黑名单
func CheckDeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	return Bean.DeleteInstanceList.Contain(*model), nil
}

// CheckLeader 检查是否是领导
func CheckLeader(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckLeader-service", util.Strval(r))
		}
	}()
	Bean.Leaders.LeaderModelsListRWLock.RLock()
	Assert := reflect.DeepEqual(*model, Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony])
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
	return Assert, nil
}
