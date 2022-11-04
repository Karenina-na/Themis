package follow

import (
	"Themis/src/entity"
	"Themis/src/exception"
	Factory "Themis/src/pool"
	"Themis/src/sync/common"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// CreateSyncAllDataRoutine
// @Description: 创建同步所有数据协程
// @param        m *syncBean.MessageModel
// @return       E error
func CreateSyncAllDataRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncAllDataRoutine-follow", util.Strval(r))
		}
	}()

	// 创建同步实例数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		err := common.DataSyncInstances(m.SyncMessage.Instances)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})

	// 创建同步黑名单数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		err := common.DataSyncDelete(m.SyncMessage.DeleteInstances)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})

	// 创建同步leader数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		err := common.DataSyncLeader(m.SyncMessage.Leaders)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}

// CreateSyncAppendInstancesDataRoutine
//
//	@Description: 创建同步实例数据协程
//	@param m	*syncBean.MessageModel
//	@return E	error
func CreateSyncAppendInstancesDataRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncAppendInstancesDataRoutine-follow", util.Strval(r))
		}
	}()

	// 创建同步实例数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		l := make([]entity.ServerModel, 0, 1)
		l = append(l, m.SyncInstance)
		err := common.DataSyncInstances(l)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}

// CreateSyncDeleteInstancesDataRoutine
//
//	@Description: 创建同步黑名单数据协程
//	@param m	*syncBean.MessageModel
//	@return E	error
func CreateSyncDeleteInstancesDataRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncDeleteInstancesDataRoutine-follow", util.Strval(r))
		}
	}()

	// 创建同步黑名单数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		l := make([]entity.ServerModel, 0, 1)
		l = append(l, m.SyncInstance)
		err := common.DataSyncDelete(l)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}

// CreateSyncCancelDeleteInstancesDataRoutine
//
//	@Description: 创建同步删除黑名单协程
//	@param m	*syncBean.MessageModel
//	@return E	error
func CreateSyncCancelDeleteInstancesDataRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncDeleteInstancesDataRoutine-follow", util.Strval(r))
		}
	}()

	// 创建同步删除黑名单协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		l := make([]entity.ServerModel, 0, 1)
		l = append(l, m.SyncInstance)
		err := common.DataSyncCancelDelete(l)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}

// CreateSyncLeaderDataRoutine
//
//	@Description: 创建同步leader数据协程
//	@param m	*syncBean.MessageModel
//	@return E	error
func CreateSyncLeaderDataRoutine(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CreateSyncLeaderDataRoutine-follow", util.Strval(r))
		}
	}()

	// 创建同步leader数据协程
	Factory.RoutinePool.CreateWork(func() (E error) {
		l := make([]entity.ServerModel, 0, 1)
		l = append(l, m.SyncInstance)
		err := common.DataSyncLeader(l)
		if err != nil {
			return err
		}
		return nil
	}, func(Message error) {
		exception.HandleException(m)
	})
	return nil
}
