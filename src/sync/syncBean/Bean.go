package syncBean

import "Themis/src/util"

var (
	//Name	该服务集群名
	Name string
	// CloseChan 关闭通道
	CloseChan chan bool
	// Term 选期
	Term int
	// Status 服务状态
	Status StatusLevel
	// LeaderName	leader名
	LeaderName string
	// LeaderAddress 选举出的leader
	LeaderAddress *SyncAddressModel
	// LeaderServicePort leader服务端口
	LeaderServicePort string
	// SyncAddress 同步服务地址
	SyncAddress *util.LinkList[SyncAddressModel]
	// UdpReceiveMessage UDP接收消息队列
	UdpReceiveMessage chan MessageModel
	// UdpSendMessage UDP发送消息队列
	UdpSendMessage chan MessageModel
)
