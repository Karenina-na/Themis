package config

import (
	"Themis/src/exception"
	"github.com/spf13/viper"
	"strconv"
)

func InitConfig() (E any) {
	defer func() {
		E = recover()
	}()
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err == nil {
		if err := LoadRoutineConfig(); err != nil {
			exception.HandleException(err)
		}
		if err := LoadPortConfig(); err != nil {
			exception.HandleException(err)
		}
		if err := LoadServerConfig(); err != nil {
			exception.HandleException(err)
		}
		if err := LoadDatabaseConfig(); err != nil {
			exception.HandleException(err)
		}
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			panic(exception.NewConfigurationPanic("init", "配置文件不存在"))
		} else {
			panic(exception.NewConfigurationPanic("init", "配置文件读取失败"+err.Error()))
		}
	}
	return nil
}

func LoadRoutineConfig() (E any) {
	defer func() {
		E = recover()
	}()
	MaxRoutineNum = viper.GetInt(`goroutine.max-goroutine`)
	if MaxRoutineNum < 1 {
		panic(exception.NewConfigurationPanic("LoadRoutineConfig", "max-goroutine非法"))
	}
	CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if CoreRoutineNum < 1 {
		panic(exception.NewConfigurationPanic("LoadRoutineConfig", "core-goroutine非法"))
	}
	if CoreRoutineNum > MaxRoutineNum {
		panic(exception.NewConfigurationPanic("LoadRoutineConfig", "core-goroutine大于max-goroutine"))
	}
	return nil
}

func LoadPortConfig() (E any) {
	defer func() {
		E = recover()
	}()
	port := viper.GetInt(`Themis.port`)
	if port < 0 || port > 65535 {
		panic(exception.NewConfigurationPanic("LoadPortConfig", "port端口非法"))
	}
	Port = strconv.Itoa(port)
	udpPort := viper.GetInt(`Themis.UDP-port`)
	if udpPort < 0 || udpPort > 65535 {
		panic(exception.NewConfigurationPanic("LoadPortConfig", "UDP-port端口非法"))
	} else if udpPort == port {
		panic(exception.NewConfigurationPanic("LoadPortConfig", "UDP-port端口不能与port端口相同"))
	}
	UDPPort = strconv.Itoa(udpPort)
	return nil
}

func LoadServerConfig() (E any) {
	defer func() {
		E = recover()
	}()
	ServerModelQueueNum = viper.GetInt(`Themis.server.model-queue`)
	if ServerModelQueueNum <= 0 {
		panic(exception.NewConfigurationPanic("LoadServerConfig", "model-queue非法"))
	}
	ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
	if ServerModelBeatQueue <= 0 {
		panic(exception.NewConfigurationPanic("LoadServerConfig", "beat-queue非法"))
	}
	ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
	if ServerBeatTime <= 0 {
		panic(exception.NewConfigurationPanic("LoadServerConfig", "beat-time非法"))
	}
	CreateLeaderAlgorithm = viper.GetString(`Themis.leader-algorithm`)
	return nil
}

func LoadDatabaseConfig() (E any) {
	defer func() {
		E = recover()
	}()
	DatabaseEnable = viper.GetBool(`Themis.database.enable`)
	if DatabaseEnable {
		PersistenceTime = int64(viper.GetInt(`Themis.database.persistence-time`))
		if PersistenceTime <= 0 {
			panic(exception.NewConfigurationPanic("LoadDatabaseConfig", "persistence-time非法"))
		}
	} else {
		PersistenceTime = 0
	}
	return nil
}
