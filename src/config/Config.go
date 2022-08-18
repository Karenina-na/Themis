package config

import (
	"Themis/src/entity/util"
	"github.com/spf13/viper"
)

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

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./src/config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err == nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			util.Loglevel(util.Info, "InitConfig", "配置文件config.yaml不存在")
		}
	}
	MaxRoutineNum = viper.GetInt(`goroutine.MaxRoutineNum`)
	CoreRoutineNum = viper.GetInt(`goroutine.CoreRoutineNum`)
	Port = viper.GetString(`Themis.Port`)
	UDPPort = viper.GetString(`Themis.UDPPort`)
	ServerModelQueueNum = viper.GetInt(`Themis.server.ServerModelQueueNum`)
	ServerModelBeatQueue = viper.GetInt(`Themis.server.ServerModelBeatQueue`)
	ServerBeatTime = int64(viper.GetInt(`Themis.server.ServerBeatTime`))
	CreateLeaderAlgorithm = viper.GetString(`Themis.CreateLeaderAlgorithm`)
}
