package leader

import (
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// StatusOperatorLeader
//
//	@Description: 领导者状态操作
//	@param m	消息
//	@param SyncRoutineBool	同步协程
//	@return B	是否关闭
//	@return E	异常
func StatusOperatorLeader(m *syncBean.MessageModel, SyncRoutineBool chan bool) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("StatusOperatorLeader-leader", util.Strval(r))
		}
	}()
	switch m.Status {
	case syncBean.LEADER:
		b, err := leaderMessageOperator(m, SyncRoutineBool)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	case syncBean.CANDIDATE:
		b, err := candidateMessageOperator(m, SyncRoutineBool)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	case syncBean.FOLLOW:
		b, err := followMessageOperator(m, SyncRoutineBool)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	}
	return false, nil
}

// leaderMessageOperator
//
//	@Description: 处理LEADER信息
//	@param m	*syncBean.MessageModel
//	@param SyncRoutineBool
//	@return B	bool
//	@return E	error
func leaderMessageOperator(m *syncBean.MessageModel, SyncRoutineBool chan bool) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("leaderMessageOperator-leader", util.Strval(r))
		}
	}()
	if m.Term > syncBean.Term {
		syncBean.Term = m.Term
		syncBean.Leader.SetLeaderModel(m.Name, m.UDPAddress.IP, m.UDPAddress.Port, m.ServicePort)
		close(SyncRoutineBool)
		util.Loglevel(util.Info, "leaderMessageOperator-leader",
			"Leader卸任,收到其他leader，成为FOLLOW")
		syncBean.Status = syncBean.FOLLOW
		return true, nil
	}
	if m.Term == syncBean.Term {
		close(SyncRoutineBool)
		util.Loglevel(util.Info, "leaderMessageOperator-leader",
			"Leader卸任,收到其他相同任期的leader，成为FOLLOW")
		syncBean.Status = syncBean.FOLLOW
		return true, nil
	}
	return false, nil
}

// candidateMessageOperator
//
//	@Description: 处理CANDIDATE信息
//	@param m	*syncBean.MessageModel
//	@param SyncRoutineBool	chan bool
//	@return B	bool
//	@return E	error
func candidateMessageOperator(m *syncBean.MessageModel, SyncRoutineBool chan bool) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("candidateMessageOperator-leader", util.Strval(r))
		}
	}()
	if m.Term > syncBean.Term {
		message := syncBean.NewMessageModel()
		message.SetMessageModeForVoteResponse(syncBean.Term, syncBean.Status,
			m.UDPAddress.IP, m.UDPAddress.Port, true)
		syncBean.UdpSendMessage <- *message
		syncBean.Term = m.Term
		close(SyncRoutineBool)
		util.Loglevel(util.Info, "Leader-sync", "Leader卸任,收到更高任期的CANDIDATE，成为FOLLOW")
		syncBean.Status = syncBean.FOLLOW
		return true, nil
	}
	return false, nil
}

// followMessageOperator
//
//	@Description: 处理FOLLOW信息
//	@param m	*syncBean.MessageModel
//	@param SyncRoutineBool	chan bool
//	@return B	bool
//	@return E	error
func followMessageOperator(m *syncBean.MessageModel, SyncRoutineBool chan bool) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("followMessageOperator-leader", util.Strval(r))
		}
	}()
	if m.Term > syncBean.Term {
		syncBean.Term = m.Term
		close(SyncRoutineBool)
		util.Loglevel(util.Info, "followMessageOperator-leader",
			"Leader卸任,收到更高任期的FOLLOW，成为FOLLOW")
		syncBean.Status = syncBean.FOLLOW
		return true, nil
	}
	return false, nil
}
