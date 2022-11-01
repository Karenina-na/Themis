package syncBean

import (
	"Themis/src/entity"
)

type MessageType int

const (
	MessageTypeHeartbeat MessageType = iota
	MessageTypeRequestVote
	MessageTypeRequestVoteResponse
	MessageTypeAppendEntries
	MessageTypeDeleteEntries
	MessageTypeCancelDeleteEntries
	MessageTypeLeaderEntries
	MessageTypeInstallSnapshot
)

// MessageModel 消息模型
type MessageModel struct {
	Type        MessageType `json:"type"` // 消息类型
	Name        string      `json:"name"`
	Term        int         `json:"term"`
	Status      StatusLevel `json:"StatusLevel"`
	SyncMessage struct {
		Instances       []entity.ServerModel `json:"instances"`
		DeleteInstances []entity.ServerModel `json:"deleteInstances"`
		Leaders         []entity.ServerModel `json:"leaders"`
	} `json:"message"`
	SyncInstance     entity.ServerModel `json:"syncInstances"`
	UDPAddress       SyncAddressModel   `json:"udpAddress"`
	UDPTargetAddress SyncAddressModel   `json:"udpTargetAddress"`
	ServicePort      string             `json:"servicePort"`
	BOOL             bool               `json:"bool"`
	Sign             string             `json:"sign"`
}
