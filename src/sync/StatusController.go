package sync

import (
	"Themis/src/exception"
	"Themis/src/sync/candidate"
	"Themis/src/sync/follow"
	"Themis/src/sync/leader"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// StatusController
//
//	@Description: 状态控制器
//	@return E	error	错误信息
func StatusController() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("StatusController-sync", util.Strval(r))
		}
	}()
	for {
		select {
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "StatusController-sync", "状态控制器退出")
			return nil
		default:
			switch syncBean.Status {
			case syncBean.FOLLOW:
				err := follow.Follow()
				if err != nil {
					exception.HandleException(err)
				}
			case syncBean.CANDIDATE:
				err := candidate.Candidate()
				if err != nil {
					exception.HandleException(err)
				}
			case syncBean.LEADER:
				err := leader.Leader()
				if err != nil {
					exception.HandleException(err)
				}
			}
		}
	}
}
