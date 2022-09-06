package syncBean

import "Themis/src/entity"

// MessageModel 消息模型
type MessageModel struct {
	Name    string      `json:"name"`
	Term    int         `json:"term"`
	Status  StatusLevel `json:"StatusLevel"`
	Message struct {
		Instances       []entity.ServerModel `json:"instances"`
		DeleteInstances []entity.ServerModel `json:"deleteInstances"`
		Leaders         []entity.ServerModel `json:"leaders"`
	} `json:"message"`
	UDPAddress       SyncAddressModel `json:"udpAddress"`
	UDPTargetAddress SyncAddressModel `json:"udpTargetAddress"`
	ServicePort      string           `json:"servicePort"`
	BOOL             bool             `json:"bool"`
}
