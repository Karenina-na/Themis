package sync

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"encoding/json"
	"net"
	"time"
)

// UDPSend 发送udp消息
func UDPSend() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("UDPSend-sync", util.Strval(r))
		}
	}()
	for {
		select {
		case msg := <-syncBean.UdpSendMessage:
			conn, err := net.DialTimeout("udp", msg.TargetAddress.IP+":"+msg.TargetAddress.Port,
				time.Duration(config.Cluster.UDPTimeOut)*time.Second)
			if err != nil {
				return exception.NewUserError("UDPSend-sync-goroutine", "UDP连接错误-"+err.Error())
			}
			data, err := json.Marshal(msg)
			if err != nil {
				return exception.NewUserError("UDPSend-sync-goroutine", "json转换错误"+err.Error())
			}
			_, err = conn.Write(data)
			if err != nil {
				return exception.NewUserError("UDPSend-sync-goroutine", "UDP发送错误-"+err.Error())
			}
			err = conn.Close()
			if err != nil {
				return exception.NewUserError("UDPSend-sync-goroutine", "UDP关闭错误-"+err.Error())
			}
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "UDPSend", "UDP发送协程退出")
			return
		}

	}
}

// UDPReceive 接收udp消息
func UDPReceive() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("UDPReceive-sync", util.Strval(r))
		}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+config.Cluster.Port)
	if err != nil {
		return exception.NewUserError("UDPReceive-sync", "创建udp服务错误"+err.Error())
	}
	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return exception.NewUserError("UDPReceive-sync", "监听udp服务错误-"+err.Error())
	}
	Bean.RoutinePool.CreateWork(func() (E error) {
		for {
			buf := make([]byte, 4096)
			n, _, err := serverConn.ReadFromUDP(buf)
			if err != nil {
				util.Loglevel(util.Debug, "UDPReceive", "UDP接收协程退出")
				return nil
			}
			var msg syncBean.MessageModel
			err = json.Unmarshal(buf[:n], &msg)
			if err != nil {
				return exception.NewUserError("UDPReceive-sync-goroutine", "json转换错误"+err.Error())
			}
			syncBean.UdpReceiveMessage <- msg
		}
	}, func(Message error) {
		exception.HandleException(Message)
	})
	<-syncBean.CloseChan
	err = serverConn.Close()
	if err != nil {
		return exception.NewUserError("UDPReceive-sync", "关闭udp服务错误-"+err.Error())
	}
	return
}
