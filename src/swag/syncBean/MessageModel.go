package syncBean

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/util/encryption"
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

// NewMessageModel
// @Description: 创建消息模型
// @return       *MessageModel 消息模型
func NewMessageModel() *MessageModel {
	return &MessageModel{
		Name: Name,
		SyncMessage: struct {
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

func (message *MessageModel) SetMessageModelForHeartBeat(
	term int, status StatusLevel,
	TargetAddressIP string, TargetAddressPort string) {
	message.Type = MessageTypeHeartbeat
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForSyncSnapShot
//
//	@Description: 设置消息模型为同步快照
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	instances	实例列表
//	@param //	deleteInstances	删除实例列表
//	@param //	leaders	领导者列表
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForSyncSnapShot( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instances []entity.ServerModel, // 实例列表
	deleteInstances []entity.ServerModel, // 删除实例列表
	leaders []entity.ServerModel, // leader列表
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeInstallSnapshot
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = instances
	message.SyncMessage.DeleteInstances = deleteInstances
	message.SyncMessage.Leaders = leaders
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForSyncAppend
//
//	@Description: 设置消息模型为同步追加
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	instances	实例列表
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForSyncAppend( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instances entity.ServerModel, // 实例列表
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeAppendEntries
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.SyncInstance = instances
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForSyncDelete
//
//	 @Description: 设置消息模型为同步删除
//	 @receiver message	消息模型
//		@receiver message	消息模型
//		@param //	term	选期
//		@param //	status	服务状态
//		@param //	instances	实例列表
//		@param //	TargetAddressIP	消息目标IP
//		@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForSyncDelete( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instance entity.ServerModel, // 实例列表
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeDeleteEntries
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.SyncInstance = instance
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForSyncCancelDelete
//
//	@Description: 设置消息模型为同步取消删除
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	instances	实例列表
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForSyncCancelDelete( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instance entity.ServerModel, // 实例列表
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeCancelDeleteEntries
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.SyncInstance = instance
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForSyncLeader
//
//	@Description: 设置消息模型为同步领导者
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	instances	实例列表
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForSyncLeader( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	instance entity.ServerModel, // 实例列表
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeLeaderEntries
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.SyncInstance = instance
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForVote
//
//	@Description: 设置消息模型为投票
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
func (message *MessageModel) SetMessageModeForVote( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
) {
	message.Type = MessageTypeRequestVote
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = false
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// SetMessageModeForVoteResponse
//
//	@Description: 设置消息模型为投票响应
//	@receiver message	消息模型
//	@param //	term	选期
//	@param //	status	服务状态
//	@param //	TargetAddressIP	消息目标IP
//	@param //	TargetAddressPort	消息目标端口库
//	@param //	B	是否同意
func (message *MessageModel) SetMessageModeForVoteResponse( // 设置消息模型
	term int, // 选期
	status StatusLevel, // 服务状态
	TargetAddressIP string, // 消息目标IP
	TargetAddressPort string, // 消息目标端口
	B bool, // 是否投票
) {
	message.Type = MessageTypeRequestVoteResponse
	message.Name = config.Cluster.ClusterName
	message.Term = term
	message.Status = status
	message.SyncMessage.Instances = nil
	message.SyncMessage.DeleteInstances = nil
	message.SyncMessage.Leaders = nil
	message.ServicePort = config.Port.CenterPort
	message.UDPAddress.IP = config.Cluster.IP
	message.UDPAddress.Port = config.Cluster.Port
	message.UDPTargetAddress.IP = TargetAddressIP
	message.UDPTargetAddress.Port = TargetAddressPort
	message.BOOL = B
	message.Sign = encryption.Sha256(message.Name +
		message.UDPAddress.IP + message.UDPAddress.Port +
		message.UDPTargetAddress.IP + message.UDPTargetAddress.Port)
}

// VerifySign
//
//	@Description: 验证消息模型签名
//	@receiver message	消息模型
//	@return bool	是否通过
func (message *MessageModel) VerifySign() bool {
	return message.Sign == encryption.Sha256(message.Name+
		message.UDPAddress.IP+message.UDPAddress.Port+
		message.UDPTargetAddress.IP+message.UDPTargetAddress.Port)
}
