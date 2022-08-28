package config

import (
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"math"
	"strconv"
)

func InitConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitConfig-config", util.Strval(r))
		}
	}()
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err == nil {
		if err := LoadRoutineConfig(); err != nil {
			return err
		}
		if err := LoadPortConfig(); err != nil {
			return err
		}
		if err := LoadServerConfig(); err != nil {
			return err
		}
		if err := LoadDatabaseConfig(); err != nil {
			return err
		}
		if err := LoadListenConfig(); err != nil {
			return err
		}
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return exception.NewConfigurationError("InitConfig-config", "配置文件不存在")
		} else {
			return exception.NewConfigurationError("InitConfig-config", "配置文件读取失败"+err.Error())
		}
	}
	return nil
}

func LoadRoutineConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadRoutineConfig-config", util.Strval(r))
		}
	}()
	MaxRoutineNum = viper.GetInt(`goroutine.max-goroutine`)
	if MaxRoutineNum < 1 {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "max-goroutine非法")
	}
	CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if CoreRoutineNum < 1 {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "core-goroutine非法")
	}
	if CoreRoutineNum > MaxRoutineNum {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "core-goroutine大于max-goroutine")
	}
	RoutineTimeOut = viper.GetInt(`goroutine.timeout`)
	if CoreRoutineNum < 1 {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "timeout非法")
	}
	return nil
}

func LoadPortConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadPortConfig-config", util.Strval(r))
		}
	}()
	port := viper.GetInt(`Themis.port`)
	if port < 0 || port > 65535 {
		return exception.NewConfigurationError("LoadPortConfig-config", "port端口非法")
	}
	Port = strconv.Itoa(port)
	udpPort := viper.GetInt(`Themis.UDP-port`)
	if udpPort < 0 || udpPort > 65535 {
		return exception.NewConfigurationError("LoadPortConfig-config", "UDP-port端口非法")
	} else if udpPort == port {
		return exception.NewConfigurationError("LoadPortConfig-config", "UDP-port端口不能与port端口相同")
	}
	UDPPort = strconv.Itoa(udpPort)
	UDPTimeOut = viper.GetInt(`Themis.UDP-timeout`)
	if UDPTimeOut < 1 {
		return exception.NewConfigurationError("LoadPortConfig-config", "UDP-timeout非法")
	}
	return nil
}

func LoadServerConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadServerConfig-config", util.Strval(r))
		}
	}()
	ServerModelQueueNum = viper.GetInt(`Themis.server.model-queue`)
	if ServerModelQueueNum <= 0 {
		return exception.NewConfigurationError("LoadServerConfig-config", "model-queue非法")
	}
	ServerModelHandleNum = viper.GetInt(`Themis.server.model-handle-number`)
	if ServerModelHandleNum <= 0 {
		return exception.NewConfigurationError("LoadServerConfig-config", "model-handle-number非法")
	}
	ServerModelBeatEnable = viper.GetBool(`Themis.server.beat-enable`)
	if ServerModelBeatEnable {
		ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
		if ServerModelBeatQueue <= 0 {
			return exception.NewConfigurationError("LoadServerConfig-config", "beat-queue非法")
		}
		ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
		if ServerBeatTime <= 0 {
			return exception.NewConfigurationError("LoadServerConfig-config", "beat-time非法")
		}
	} else {
		ServerModelBeatQueue = 0
		ServerBeatTime = math.MaxInt
	}
	CreateLeaderAlgorithm = viper.GetString(`Themis.leader-algorithm`)
	return nil
}

func LoadDatabaseConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadDatabaseConfig-config", util.Strval(r))
		}
	}()
	DatabaseEnable = viper.GetBool(`Themis.database.enable`)
	if DatabaseEnable {
		PersistenceTime = int64(viper.GetInt(`Themis.database.persistence-time`))
		if PersistenceTime <= 0 {
			return exception.NewConfigurationError("LoadDatabaseConfig-config", "persistence-time非法")
		}
	} else {
		PersistenceTime = 0
	}
	return nil
}

func LoadListenConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadListenConfig-config", util.Strval(r))
		}
	}()
	ListenTime = int64(viper.GetInt(`Themis.listen.space-time`))
	if ListenTime <= 0 {
		return exception.NewConfigurationError("LoadDatabaseConfig-config", "space-time非法")
	}
	return nil
}
