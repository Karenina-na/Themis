package config

var (
	// MaxRoutineNum goroutine池最大线程数
	MaxRoutineNum int
	// CoreRoutineNum goroutine池核心线程数
	CoreRoutineNum int

	// Port 注册中心http端口
	Port string

	// UDPPort UDP服务端口
	UDPPort string

	// ServerModelQueueNum 服务注册处理队列容量
	ServerModelQueueNum int

	// ServerModelBeatQueue 服务心跳处理队列容量
	ServerModelBeatQueue int

	// ServerBeatTime 服务心跳超时时间   单位：s
	ServerBeatTime int64

	// CreateLeaderAlgorithm 记账人选举算法
	CreateLeaderAlgorithm string
)
