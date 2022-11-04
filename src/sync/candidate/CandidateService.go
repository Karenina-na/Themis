package candidate

import (
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// BroadcastForVote
//
//	@Description: BroadcastForVote
//	@return E	: error
func BroadcastForVote() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("BroadcastForVote-candidate", util.Strval(r))
		}
	}()

	// 广播投票
	syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
		m := syncBean.NewMessageModel()
		m.SetMessageModeForVote(syncBean.Term, syncBean.Status, value.IP, value.Port)
		syncBean.UdpSendMessage <- *m
	})
	return nil
}

// JudgeVoteEnough
//
//	@Description: JudgeVoteEnough
//	@param m	: syncBean.MessageModel
//	@param num	: int
//	@param AllNum	: int
//	@return B	: bool
//	@return E	: error
func JudgeVoteEnough(num *int, AllNum *int) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("JudgeVoteEnough-candidate", util.Strval(r))
		}
	}()

	// 判断是否获得了大多数的投票
	if *num > *AllNum/2 {
		util.Loglevel(util.Debug, "Candidate-candidate", "选票足够，转换LEADER")
		syncBean.Status = syncBean.LEADER
		return true, nil
	}
	return false, nil
}

// StatusOperatorCandidate
//
//	@Description: StatusOperatorCandidate
//	@param num	: int
//	@param AllNum	: int
//	@return sign	: bool
//	@return E	: error
func StatusOperatorCandidate(m *syncBean.MessageModel, num *int, AllNum *int) (sign bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("StatusOperatorCandidate-candidate", util.Strval(r))
		}
	}()
	//根据收到的消息发送者状态进行不同的操作，返回是否需要继续等待
	switch m.Status {
	case syncBean.LEADER:
		b, err := leaderMessageOperator(m)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	case syncBean.CANDIDATE:
		b, err := candidateMessageOperator(m)
		if err != nil {
			return false, err
		}
		if b {
			return true, nil
		}
	}

	//判断任期，添加选票
	if m.Term < syncBean.Term {
		if m.BOOL {
			*num++
		}
	}

	//判断是否获得了大多数的投票
	b, E := JudgeVoteEnough(num, AllNum)
	if E != nil {
		return false, E
	}
	if b {
		return true, nil
	}
	return false, nil
}

// leaderMessageOperator
//
//	@Description: leaderMessageOperator
//	@param m	: syncBean.MessageModel
//	@return B	: bool
//	@return E	: error
func leaderMessageOperator(m *syncBean.MessageModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("leaderMessageOperator-candidate", util.Strval(r))
		}
	}()

	//判断数据类型
	switch m.Type {
	case syncBean.MessageTypeInstallSnapshot, syncBean.MessageTypeHeartbeat,
		syncBean.MessageTypeLeaderEntries, syncBean.MessageTypeDeleteEntries, syncBean.MessageTypeCancelDeleteEntries,
		syncBean.MessageTypeAppendEntries:

		//如果任期大于等于自己的任期，转换FOLLOW
		if m.Term >= syncBean.Term {
			syncBean.Term = m.Term
			syncBean.Leader.SetLeaderModel(m.Name, m.UDPAddress.IP, m.UDPAddress.Port, m.ServicePort)
			util.Loglevel(util.Debug, "Candidate-candidate", "收到Leader，转换FOLLOW")
			syncBean.Status = syncBean.FOLLOW
			return true, nil
		}
	}
	return false, nil
}

// candidateMessageOperator
//
//	@Description: leaderMessageOperator
//	@param m	: syncBean.MessageModel
//	@return B	: bool
//	@return E	: error
func candidateMessageOperator(m *syncBean.MessageModel) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("candidateMessageOperator-candidate", util.Strval(r))
		}
	}()

	//如果任期大于等于自己的任期，转换FOLLOW，并投票给他
	if m.Term > syncBean.Term {
		message := syncBean.NewMessageModel()
		message.SetMessageModeForVoteResponse(syncBean.Term, syncBean.Status,
			m.UDPAddress.IP, m.UDPAddress.Port, true)
		syncBean.UdpSendMessage <- *message
		util.Loglevel(util.Debug, "Candidate-candidate", "收到更大Term的Candidate，转换FOLLOW")
		syncBean.Status = syncBean.FOLLOW
		return true, nil
	}
	return false, nil
}
