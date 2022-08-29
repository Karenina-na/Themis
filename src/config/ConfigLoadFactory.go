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
	if !VerifyReg(PositiveReg, strconv.Itoa(MaxRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.max-goroutine非法")
	}
	CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if !VerifyReg(PositiveReg, strconv.Itoa(CoreRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.core-goroutine非法")
	}
	if CoreRoutineNum > MaxRoutineNum {
		return exception.NewConfigurationError("LoadRoutineConfig-config",
			"goroutine.core-goroutine大于goroutine.max-goroutine")
	}
	RoutineTimeOut = viper.GetInt(`goroutine.timeout`)
	if !VerifyReg(PositiveReg, strconv.Itoa(RoutineTimeOut)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.timeout非法")
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
	Port = viper.GetString(`Themis.port`)
	if !VerifyReg(PortReg, Port) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.port端口非法")
	}
	UDPPort = viper.GetString(`Themis.UDP-port`)
	if !VerifyReg(PortReg, UDPPort) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.UDP-port端口非法")
	} else if UDPPort == Port {
		return exception.NewConfigurationError("LoadPortConfig-config",
			"Themis.UDP-port端口不能与Themis.port端口相同")
	}
	udpTimeOut := viper.GetString(`Themis.UDP-timeout`)
	if !VerifyReg(PositiveReg, udpTimeOut) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.UDP-timeout非法")
	}
	UDPTimeOut, _ = strconv.Atoi(udpTimeOut)
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
	if !VerifyReg(PositiveReg, strconv.Itoa(ServerModelQueueNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-queue非法")
	}
	ServerModelHandleNum = viper.GetInt(`Themis.server.model-handle-number`)
	if !VerifyReg(PositiveReg, strconv.Itoa(ServerModelHandleNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-handle-number")
	}
	ServerModelBeatEnable = viper.GetBool(`Themis.server.beat-enable`)
	if ServerModelBeatEnable {
		ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
		if !VerifyReg(PositiveReg, strconv.Itoa(ServerModelBeatQueue)) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-queue非法")
		}
		ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(ServerBeatTime))) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-time非法")
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
		if !VerifyReg(PositiveReg, strconv.Itoa(int(PersistenceTime))) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config", "Themis.database.persistence-time非法")
		}
		DatabaseSoftDeleteEnable = viper.GetBool(`Themis.database.soft-delete-enable`)
		Database.DatabaseType = viper.GetString(`Themis.database.type`)
		Database.DatabaseHost = viper.GetString(`Themis.database.mysql.host`)
		if !VerifyReg(IpReg, Database.DatabaseHost) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config",
				"Themis.database.mysql.host非法")
		}
		Database.DatabasePort = viper.GetString(`Themis.database.mysql.port`)
		if !VerifyReg(PortReg, Database.DatabasePort) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config",
				"Themis.database.mysql.port非法")
		}
		Database.DatabaseName = viper.GetString(`Themis.database.mysql.name`)
		Database.DatabaseUser = viper.GetString(`Themis.database.mysql.user`)
		Database.DatabasePassword = viper.GetString(`Themis.database.mysql.password`)
		Database.DatabaseMaxOpenConns = viper.GetInt(`Themis.database.mysql.max-open-conns`)
		if !VerifyReg(PositiveReg, strconv.Itoa(Database.DatabaseMaxOpenConns)) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config",
				"Themis.database.mysql.max-open-conns非法")
		}
		Database.DatabaseMaxIdleConns = viper.GetInt(`Themis.database.mysql.max-idle-conns`)
		if !VerifyReg(PositiveReg, strconv.Itoa(Database.DatabaseMaxIdleConns)) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config",
				"Themis.database.mysql.max-idle-conns非法")
		}
		Database.DatabaseMaxLifetimeConns = viper.GetInt(`Themis.database.mysql.max-conns-lifetime`)
		if !VerifyReg(PositiveReg, strconv.Itoa(Database.DatabaseMaxLifetimeConns)) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config",
				"Themis.database.mysql.max-conns-lifetime非法")
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
	if !VerifyReg(PositiveReg, strconv.Itoa(int(ListenTime))) {
		return exception.NewConfigurationError("LoadListenConfig-config", "Themis.listen.space-time非法")
	}
	return nil
}
