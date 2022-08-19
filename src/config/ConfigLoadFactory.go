package config

import (
	"Themis/src/entity/util"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err == nil {
		LoadRoutineConfig()
		LoadPortConfig()
		LoadServerConfig()
		LoadDatabaseConfig()
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			util.Loglevel(util.Info, "InitConfig", "配置文件config.yaml不存在-"+err.Error())
		} else {
			util.Loglevel(util.Warn, "InitConfig", "未知错误-"+err.Error())
		}
		os.Exit(0)
	}
}

func LoadRoutineConfig() {
	MaxRoutineNum = viper.GetInt(`goroutine.max-goroutine`)
	if MaxRoutineNum < 1 {
		util.Loglevel(util.Info, "InitConfig", "max-goroutine")
		os.Exit(0)
	}
	CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if CoreRoutineNum < 1 {
		util.Loglevel(util.Info, "InitConfig", "core-goroutine")
		os.Exit(0)
	}
	if CoreRoutineNum > MaxRoutineNum {
		util.Loglevel(util.Info, "InitConfig", "max-goroutine不能比core-goroutine小")
		os.Exit(0)
	}
}

func LoadPortConfig() {
	port := viper.GetInt(`Themis.port`)
	if port < 0 && port > 65535 {
		util.Loglevel(util.Info, "InitConfig", "port端口错误")
		os.Exit(0)
	}
	Port = strconv.Itoa(port)
	udpPort := viper.GetInt(`Themis.UDP-port`)
	if udpPort < 0 && udpPort > 65535 {
		util.Loglevel(util.Info, "InitConfig", "UDP-port端口错误")
		os.Exit(0)
	} else if udpPort == port {
		util.Loglevel(util.Info, "InitConfig", "port与UDP-port端口冲突")
		os.Exit(0)
	}
	UDPPort = strconv.Itoa(udpPort)
}

func LoadServerConfig() {
	ServerModelQueueNum = viper.GetInt(`Themis.server.model-queue`)
	if ServerModelQueueNum <= 0 {
		util.Loglevel(util.Info, "InitConfig", "model-queue非法")
		os.Exit(0)
	}
	ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
	if ServerModelBeatQueue <= 0 {
		util.Loglevel(util.Info, "InitConfig", "beat-queue非法")
		os.Exit(0)
	}
	ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
	if ServerModelBeatQueue <= 0 {
		util.Loglevel(util.Info, "InitConfig", "beat-time非法")
		os.Exit(0)
	}
	CreateLeaderAlgorithm = viper.GetString(`Themis.leader-algorithm`)
}

func LoadDatabaseConfig() {
	DatabaseEnable = viper.GetBool(`Themis.database.enable`)
	if DatabaseEnable {
		PersistenceTime = int64(viper.GetInt(`Themis.database.persistence-time`))
		if PersistenceTime <= 0 {
			util.Loglevel(util.Info, "InitConfig", "persistence-time非法")
			os.Exit(0)
		}
	} else {
		PersistenceTime = 0
	}
}
