package sync

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"math/rand"
	"time"
)

// Candidate 候选人状态
func Candidate() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Candidate-sync", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Candidate-sync", "CANDIDATE状态")
CandidateHead:
	syncBean.Status = syncBean.CANDIDATE
	syncBean.Term++
	syncBean.SyncAddress.Iterator(func(index int, value syncBean.SyncAddressModel) {
		m := syncBean.NewMessageModel()
		m.SetMessageMode(syncBean.Term, syncBean.Status,
			nil, nil, nil,
			config.Cluster.IP, config.Cluster.Port,
			value.IP, value.Port, false)
		syncBean.UdpSendMessage <- *m
	})
	var num int
	num++
	AllNum := len(config.Cluster.Clusters) + 1
	for {
		select {
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "Candidate-sync", "CANDIDATE协程退出")
			return nil
		case m := <-syncBean.UdpReceiveMessage:
			switch m.Status {
			case syncBean.LEADER:
				if m.Term >= syncBean.Term {
					syncBean.Term = m.Term
					syncBean.LeaderAddress.IP = m.Address.IP
					syncBean.LeaderAddress.Port = m.Address.Port
					util.Loglevel(util.Debug, "Candidate-sync", "收到Leader，转换FOLLOW")
					E := ChangeToFollow()
					if E != nil {
						return E
					}
					return nil
				}
				if m.Term < syncBean.Term {
					if m.BOOL {
						num++
						if num >= AllNum/2 {
							util.Loglevel(util.Debug, "Candidate-sync", "收到大多数Candidate，转换FOLLOW")
							E := ChangeToFollow()
							if E != nil {
								return E
							}
							return nil
						}
					}
				}
			case syncBean.CANDIDATE:
				if m.Term > syncBean.Term {
					message := syncBean.NewMessageModel()
					message.SetMessageMode(syncBean.Term, syncBean.Status,
						nil, nil, nil,
						config.Cluster.IP, config.Cluster.Port,
						m.Address.IP, m.Address.Port, true)
					syncBean.UdpSendMessage <- *message
					util.Loglevel(util.Debug, "Candidate-sync", "收到更大Term的Candidate，转换FOLLOW")
					E := ChangeToFollow()
					if E != nil {
						return E
					}
					return nil
				}
				if m.Term < syncBean.Term {
					if m.BOOL {
						num++
						if num > AllNum/2 {
							util.Loglevel(util.Debug, "Candidate-sync", "选票足够，转换LEADER")
							E := ChangeToLeader()
							if E != nil {
								return E
							}
							return nil
						}
					}
				}
			case syncBean.FOLLOW:
				if m.BOOL {
					num++
					if num > AllNum/2 {
						util.Loglevel(util.Debug, "Candidate-sync", "选票足够，转换LEADER")
						E := ChangeToLeader()
						if E != nil {
							return E
						}
						return nil
					}
				}
			}
		case <-time.After(time.Second * time.Duration(rand.Int()%
			int(config.Cluster.MaxCandidateTimeOut-config.Cluster.MinCandidateTimeOut)+int(config.Cluster.MinCandidateTimeOut))):
			if num > AllNum/2 {
				util.Loglevel(util.Debug, "Candidate-sync", "选举超时，选票足够，转换LEADER")
				E := ChangeToLeader()
				if E != nil {
					return E
				}
				return nil
			}
			util.Loglevel(util.Debug, "Candidate-sync", "选举超时，选票不足，重新选票")
			goto CandidateHead
		}
	}
}
