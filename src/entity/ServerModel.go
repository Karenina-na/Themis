package entity

import (
	"encoding/json"
	"net"
	"time"
)

type ServerModel struct {
	IP        string `json:"IP" gorm:"column:ip"`
	Port      string `json:"port" gorm:"column:port"`
	Name      string `json:"name" gorm:"column:name"`
	Time      string `json:"time" gorm:"column:time"`
	Colony    string `json:"colony" gorm:"column:colony"`
	Namespace string `json:"namespace" gorm:"column:namespace"`
}

func NewServerModel() *ServerModel {
	return &ServerModel{}
}

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
