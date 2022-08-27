package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
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
	return InstanceList.Contain(*model), nil
}

// CheckDeleteServer 检查服务是否存在黑名单
func CheckDeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	return DeleteInstanceList.Contain(*model), nil
}

// CheckLeader 检查是否是领导
func CheckLeader(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckLeader-service", util.Strval(r))
		}
	}()
	LeadersRWLock.RLock()
	Assert := reflect.DeepEqual(*model, Leaders[model.Namespace][model.Colony])
	LeadersRWLock.RUnlock()
	return Assert, nil
}
