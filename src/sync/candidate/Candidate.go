package candidate

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"math/rand"
	"time"
)

// Candidate
// @Description: Candidate Service
// @return       E : error
func Candidate() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Candidate-candidate", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Candidate-candidate", "CANDIDATE状态")

	//任期加一
	syncBean.Term++

	//广播拉取选票
	if err := BroadcastForVote(); err != nil {
		return err
	}

	//选票
	var num int
	num++
	AllNum := len(config.Cluster.Clusters) + 1
	for {
		select {
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "Candidate-candidate", "CANDIDATE退出")
			return nil
		case m := <-syncBean.UdpReceiveMessage:

			//处理选票
			sign, E := StatusOperatorCandidate(&m, &num, &AllNum)
			if E != nil {
				return E
			}
			if sign {
				return nil
			}
		case <-time.After(time.Second * time.Duration(rand.Int()%
			int(config.Cluster.MaxCandidateTimeOut-config.Cluster.MinCandidateTimeOut)+
			int(config.Cluster.MinCandidateTimeOut))):

			//超时，判断选票是否足够
			b, E := JudgeVoteEnough(&num, &AllNum)
			if E != nil {
				return E
			}
			if b {
				return nil
			}
			util.Loglevel(util.Debug, "Candidate-candidate", "选举超时，选票不足，重新选票")
			return nil
		}
	}
}
