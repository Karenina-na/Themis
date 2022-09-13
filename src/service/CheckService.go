package service

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"reflect"
)

// CheckServer
// @Description: 检查服务是否存在
// @param        model 服务模型
// @return       B     是否存在
// @return       E     错误
func CheckServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckServer-service", util.Strval(r))
		}
	}()
	return Bean.InstanceList.Contain(*model), nil
}

// CheckDeleteServer
// @Description: 检查服务是否在黑名单
// @param        model 服务模型
// @return       B     是否存在
// @return       E     错误
func CheckDeleteServer(model *entity.ServerModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("CheckDeleteServer-service", util.Strval(r))
		}
	}()
	return Bean.DeleteInstanceList.Contain(*model), nil
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
	Assert := reflect.DeepEqual(*model, Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony])
	Bean.Leaders.LeaderModelsListRWLock.RUnlock()
	return Assert, nil
}
