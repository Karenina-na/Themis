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
	Name      string `json:"name" gorm:"column:name"`
	Time      string `json:"time" gorm:"column:time"`
	Colony    string `json:"colony" gorm:"column:colony"`
	Namespace string `json:"namespace" gorm:"column:namespace"`
}

//
// NewServerModel
// @Description: Create a new server model
// @return       *ServerModel : The new server model
//
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

func (server ServerModel) SendMessageUDP(leader ServerModel, port string, timeout int) error {
	conn, err := net.DialTimeout("udp", server.IP+":"+port,
		time.Duration(timeout)*time.Second)
	if err != nil {
		return err
	} else {
		data, _ := json.Marshal(leader)
		_, _ = conn.Write(data)
	}
	return nil
}
