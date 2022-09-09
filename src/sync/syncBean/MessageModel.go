package syncBean

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/util/encryption"
)

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
	Sign             string           `json:"sign"`
}

//
// NewMessageModel
// @Description: 创建消息模型
// @return       *MessageModel 消息模型
//
func NewMessageModel() *MessageModel {
	return &MessageModel{
		Name: Name,
		Message: struct {
			Instances       []entity.ServerModel `json:"instances"`
			DeleteInstances []entity.ServerModel `json:"deleteInstances"`
			Leaders         []entity.ServerModel `json:"leaders"`
		}(struct {
			Instances       []entity.ServerModel
			DeleteInstances []entity.ServerModel
			Leaders         []entity.ServerModel
		}{Instances: make([]entity.ServerModel, 0),
			DeleteInstances: make([]entity.ServerModel, 0),
			Leaders:         make([]entity.ServerModel, 0)}),
		UDPAddress:       SyncAddressModel{},
		UDPTargetAddress: SyncAddressModel{},
	}
}

//
// SetMessageMode
// @Description: 设置消息模型
// @receiver     message 消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //服务地址 消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
// @param        //     消息模型
//
func (message *MessageModel) SetMessageMode( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instances []entity.ServerModel, // 实例列表
	deleteInstances []entity.ServerModel, // 删除实例列表
	leaders []entity.ServerModel, // leader列表
	servicePort string, //服务地址
	AddressIP string, // 消息来源IP
	AddressPort string, // 消息来源端口
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
	BOOL bool, // 布尔值
) {
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.Message.Instances = instances
	message.Message.DeleteInstances = deleteInstances
	message.Message.Leaders = leaders
	message.ServicePort = servicePort
	message.UDPAddress.IP = AddressIP
	message.UDPAddress.Port = AddressPort
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = BOOL
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

//
// VerifySign
//  @Description: 验证消息模型签名
//  @receiver message	消息模型
//  @return bool	是否通过
//
func (message *MessageModel) VerifySign() bool {
	return message.Sign == encryption.Sha256(message.Name+
		message.UDPAddress.IP+message.UDPAddress.Port+
		message.UDPTargetAddress.IP+message.UDPTargetAddress.Port)
}