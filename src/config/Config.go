package config

var (
	// MaxRoutineNum goroutine池最大线程数
	MaxRoutineNum int
	// CoreRoutineNum goroutine池核心线程数
	CoreRoutineNum int
	// RoutineTimeOut goroutine池线程超时时间
	RoutineTimeOut int

	// Port 注册中心http端口
	Port string

	// UDPPort UDP服务端口
	UDPPort string
	// UDPTimeOut UDP超时时间
	UDPTimeOut int

	// ServerModelQueueNum 服务注册处理队列容量
	ServerModelQueueNum int

	// ServerModelHandleNum 服务注册处理器数量
	ServerModelHandleNum int

	// ServerModelBeatEnable 服务注册心跳开关
	ServerModelBeatEnable bool

	// ServerModelBeatQueue 服务心跳处理队列容量
	ServerModelBeatQueue int

	// ServerBeatTime 服务心跳超时时间   单位：s
	ServerBeatTime int64

	// CreateLeaderAlgorithm 记账人选举算法
	CreateLeaderAlgorithm string

	// Database 数据库模型
	Database struct {
		// DatabaseType 数据库类型
		DatabaseType string

		// DatabaseHost 数据库地址
		DatabaseHost string

		// DatabasePort 数据库端口
		DatabasePort string

		// DatabaseName 数据库名称
		DatabaseName string

		// DatabaseUser 数据库用户名
		DatabaseUser string

		// DatabasePassword 数据库密码
		DatabasePassword string

		// DatabaseMaxOpenConns 数据库最大连接数
		DatabaseMaxOpenConns int

		// DatabaseMaxIdleConns 数据库最大空闲连接数
		DatabaseMaxIdleConns int

		// DatabaseMaxLifetimeConns 数据库最大连接生命周期
		DatabaseMaxLifetimeConns int
	}

	// DatabaseEnable 是否开启持久化
	DatabaseEnable bool

	// PersistenceTime 持久化时间
	PersistenceTime int64

	// DatabaseSoftDeleteEnable 是否开启软删除
	DatabaseSoftDeleteEnable bool

	// ListenTime 监听时间
	ListenTime int64
)
