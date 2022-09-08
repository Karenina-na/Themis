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
	// Leader 选举出的Leader模型
	Leader LeaderModel
	// SyncAddress 同步服务地址
	SyncAddress *util.LinkList[SyncAddressModel]
	// UdpReceiveMessage UDP接收消息队列
	UdpReceiveMessage chan MessageModel
	// UdpSendMessage UDP发送消息队列
	UdpSendMessage chan MessageModel
)
