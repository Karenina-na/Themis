package follow

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"math/rand"
	"time"
)

// Follow
// @Description: 跟随者
// @return       E error
func Follow() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("Follow-follow", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Follow-follow", "FOLLOW状态")
	for {
		select {
		case <-time.After(time.Second * time.Duration(rand.Int()%
			int(config.Cluster.MaxFollowTimeOut-config.Cluster.MinFollowTimeOut)+int(config.Cluster.MinFollowTimeOut))):
			util.Loglevel(util.Debug, "Follow-follow", "FOLLOW超时，成为CANDIDATE")
			syncBean.Status = syncBean.CANDIDATE
			return nil
		case m := <-syncBean.UdpReceiveMessage:
			if err := StatusOperatorFollow(&m); err != nil {
				return err
			}
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "Follow-follow", "FOLLOW退出")
			return nil
		}
	}
}
