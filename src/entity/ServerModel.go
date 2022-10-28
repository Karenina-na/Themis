package entity

import (
	"encoding/json"
	"net"
	"time"
)

// ServerModel is a struct that contains the information of a server
type ServerModel struct {
	IP        string `json:"IP" gorm:"column:ip"`
	Port      string `json:"port" gorm:"column:port"`
	UDPPort   string `json:"udp_port" gorm:"column:udp_port"`
	Name      string `json:"name" gorm:"column:name"`
	Time      string `json:"time" gorm:"column:time"`
	Colony    string `json:"colony" gorm:"column:colony"`
	Namespace string `json:"namespace" gorm:"column:namespace"`
}

// NewServerModel
// @Description: Create a new server model
// @return       *ServerModel : The new server model
func NewServerModel() *ServerModel {
	return &ServerModel{}
}

//
// Clone
// @Description: Clone the server model
// @receiver     server       : The server model to clone
// @return       *ServerModel : The cloned server model
//

func (server ServerModel) Clone() *ServerModel {
	return &ServerModel{
		IP:        server.IP,
		Port:      server.Port,
		UDPPort:   server.UDPPort,
		Name:      server.Name,
		Time:      server.Time,
		Colony:    server.Colony,
		Namespace: server.Namespace,
	}
}

//
// SendMessageUDP
// @Description: Send a message to the server
// @receiver     server  : The server model
// @param        leader  : The leader  server model
// @param        port    : The port    to     send the message
// @param        timeout : The timeout to     send the message
// @return       error   : The error if there is one
//

func (server ServerModel) SendMessageUDP(leader ServerModel, timeout int) error {
	conn, err := net.DialTimeout("udp", server.IP+":"+server.UDPPort,
		time.Duration(timeout)*time.Second)
	if err != nil {
		return err
	} else {
		data, _ := json.Marshal(leader)
		_, _ = conn.Write(data)
	}
	return nil
}

// Equal
//
//	@Description: Check if the server model is equal to another server model
//	@receiver serverModel
//	@param s
//	@return bool
func (serverModel *ServerModel) Equal(s *ServerModel) bool {
	if serverModel.Name != s.Name {
		return false
	}
	if serverModel.IP != s.IP {
		return false
	}
	if serverModel.Port != s.Port {
		return false
	}
	if serverModel.UDPPort != s.UDPPort {
		return false
	}
	if serverModel.Colony != s.Colony {
		return false
	}
	if serverModel.Namespace != s.Namespace {
		return false
	}
	return true
}
