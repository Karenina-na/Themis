package follow

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
)

// StatusOperatorFollow
// @Description: FOLLOW状态下的操作
// @param        m *syncBean.MessageModel
// @return       E error
func StatusOperatorFollow(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("StatusOperatorFollow-follow", util.Strval(r))
		}
	}()
	switch m.Status {
	case syncBean.LEADER:
		if m.Term >= syncBean.Term {
			if err := leaderMessageOperator(m); err != nil {
				return err
			}
		}
	case syncBean.CANDIDATE:
		if m.Term > syncBean.Term && m.Type == syncBean.MessageTypeRequestVote {
			if err := candidateMessageOperator(m); err != nil {
				return err
			}
		}
	}
	return nil
}

// leaderMessageOperator
//
//	@Description: 处理LEADER信息
//	@param m	*syncBean.MessageModel
//	@return E	error
func leaderMessageOperator(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("leaderMessageOperator-follow", util.Strval(r))
		}
	}()
	if m.UDPAddress.IP != syncBean.Leader.LeaderAddress.IP ||
		m.UDPAddress.Port != syncBean.Leader.LeaderAddress.Port {
		syncBean.Leader.SetLeaderModel(m.Name, m.UDPAddress.IP, m.UDPAddress.Port, m.ServicePort)
		syncBean.Term = m.Term
	}
	switch m.Type {
	case syncBean.MessageTypeInstallSnapshot:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeInstallSnapshot-"+util.Strval(m.UDPAddress))
		}
		if err := CreateSyncAllDataRoutine(m); err != nil {
			return err
		}
	case syncBean.MessageTypeAppendEntries:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeAppendEntries-"+util.Strval(m.UDPAddress))
		}
		if err := CreateSyncAppendInstancesDataRoutine(m); err != nil {
			return err
		}
	case syncBean.MessageTypeDeleteEntries:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeDeleteEntries-"+util.Strval(m.UDPAddress))
		}
		if err := CreateSyncDeleteInstancesDataRoutine(m); err != nil {
			return err
		}
	case syncBean.MessageTypeCancelDeleteEntries:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeCancelDeleteEntries-"+util.Strval(m.UDPAddress))
		}
		if err := CreateSyncCancelDeleteInstancesDataRoutine(m); err != nil {
			return err
		}
	case syncBean.MessageTypeLeaderEntries:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeCancelDeleteEntries-"+util.Strval(m.UDPAddress))
		}
		if err := CreateSyncLeaderDataRoutine(m); err != nil {
			return err
		}
	case syncBean.MessageTypeHeartbeat:
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "StatusOperatorFollow-follow",
				"收到LEADER信息-MessageTypeHeartbeat-"+util.Strval(m.UDPAddress))
		}
	}
	return nil
}

// candidateMessageOperator
//
//	@Description: 处理CANDIDATE信息
//	@param m	*syncBean.MessageModel
//	@return E	error
func candidateMessageOperator(m *syncBean.MessageModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("candidateMessageOperator-follow", util.Strval(r))
		}
	}()
	if config.Cluster.TrackEnable {
		util.Loglevel(util.Debug, "candidateMessageOperator-follow", "收到CANDIDATE信息-"+util.Strval(m.UDPAddress))
	}
	if m.Type == syncBean.MessageTypeRequestVote {
		message := syncBean.NewMessageModel()
		message.SetMessageModeForVoteResponse(syncBean.Term, syncBean.Status,
			m.UDPAddress.IP, m.UDPAddress.Port, true)
		if config.Cluster.TrackEnable {
			util.Loglevel(util.Debug, "candidateMessageOperator-follow", "投票true")
		}
		syncBean.UdpSendMessage <- *message
	}
	return nil
}
