package common

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/service/Bean"
	"Themis/src/sync/syncBean"
	"Themis/src/util"
	"Themis/src/util/encryption"
	"encoding/json"
	"net"
	"time"
)

// UDPSend
// @Description: 发送udp消息
// @return       E error
func UDPSend() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("UDPSend-common", util.Strval(r))
		}
	}()
	for {
		select {
		case msg := <-syncBean.UdpSendMessage:
			msg.ServicePort = config.Port.CenterPort
			conn, err := net.DialTimeout("udp", msg.UDPTargetAddress.IP+":"+msg.UDPTargetAddress.Port,
				time.Duration(config.Cluster.UDPTimeOut)*time.Second)
			if err != nil {
				return exception.NewUserError("UDPSend-common", "UDP连接错误-"+err.Error())
			}
			data, err := json.Marshal(msg)
			if err != nil {
				return exception.NewUserError("UDPSend-common", "json转换错误"+err.Error())
			}
			if config.Cluster.EnableEncryption {
				data = []byte(encryption.AESEncrypt(string(data), config.Cluster.EncryptionKey))
			}
			_, err = conn.Write(data)
			if err != nil {
				return exception.NewUserError("UDPSend-common", "UDP发送错误-"+err.Error())
			}
			err = conn.Close()
			if err != nil {
				return exception.NewUserError("UDPSend-common", "UDP关闭错误-"+err.Error())
			}
		case <-syncBean.CloseChan:
			util.Loglevel(util.Debug, "UDPSend", "UDP发送协程退出")
			return
		}

	}
}

// UDPReceive
// @Description: 接收udp消息
// @return       E error
func UDPReceive() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("UDPReceive-common", util.Strval(r))
		}
	}()
	udpAddr, err := net.ResolveUDPAddr("udp", ":"+config.Cluster.Port)
	if err != nil {
		return exception.NewUserError("UDPReceive-common", "创建udp服务错误"+err.Error())
	}
	serverConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return exception.NewUserError("UDPReceive-common", "监听udp服务错误-"+err.Error())
	}
	Bean.RoutinePool.CreateWork(func() (E error) {
		for {
			buf := make([]byte, 4096)
			n, _, err := serverConn.ReadFromUDP(buf)
			if err != nil {
				util.Loglevel(util.Debug, "UDPReceive-common", "UDP接收协程退出")
				return nil
			}
			var msg syncBean.MessageModel
			if config.Cluster.EnableEncryption {
				buf = []byte(encryption.AESDecrypt(string(buf[:n]), config.Cluster.EncryptionKey))
			} else {
				buf = buf[:n]
			}
			err = json.Unmarshal(buf, &msg)
			if err != nil {
				return exception.NewUserError("UDPReceive-common", "json转换错误"+err.Error())
			}
			if msg.VerifySign() {
				syncBean.UdpReceiveMessage <- msg
			} else {
				util.Loglevel(util.Info, "UDPReceive-common", "签名错误")
			}
		}
	}, func(Message error) {
		exception.HandleException(Message)
	})
	<-syncBean.CloseChan
	err = serverConn.Close()
	if err != nil {
		return exception.NewUserError("UDPReceive-common", "关闭udp服务错误-"+err.Error())
	}
	return
}
