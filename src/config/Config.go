package config

var (
	// MaxRoutineNum goroutine池最大线程数
	MaxRoutineNum = 2000
	// CoreRoutineNum goroutine池核心线程数
	CoreRoutineNum = 1000

	// Port 注册中心http端口
	Port = "8088"

	// UDPPort UDP服务端口
	UDPPort = "12345"

	// ServerModelQueueNum 服务注册处理队列容量
	ServerModelQueueNum = 100

	// ServerModelBeatQueue 服务心跳处理队列容量
	ServerModelBeatQueue = 100

	// ServerBeatTime 服务心跳超时时间   单位：s
	ServerBeatTime int64 = 5000

	// CreateLeaderAlgorithm 记账人选举算法
	CreateLeaderAlgorithm = "RandomAlgorithmCreateLeader"
)
