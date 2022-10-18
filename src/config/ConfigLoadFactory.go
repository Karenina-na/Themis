package config

import (
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"math"
	"strconv"
)

// InitConfig
// @Description: 初始化配置文件
// @return       E error
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
		if err := LoadClusterConfig(); err != nil {
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

// LoadRoutineConfig
// @Description: 加载协程配置
// @return       E error
func LoadRoutineConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadRoutineConfig-config", util.Strval(r))
		}
	}()
	Goroutine.MaxRoutineNum = viper.GetInt(`goroutine.max-goroutine`)
	if !VerifyReg(PositiveReg, strconv.Itoa(Goroutine.MaxRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.max-goroutine非法")
	}
	Goroutine.CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if !VerifyReg(PositiveReg, strconv.Itoa(Goroutine.CoreRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.core-goroutine非法")
	}
	if Goroutine.CoreRoutineNum > Goroutine.MaxRoutineNum {
		return exception.NewConfigurationError("LoadRoutineConfig-config",
			"goroutine.core-goroutine大于goroutine.max-goroutine")
	}
	Goroutine.RoutineTimeOut = viper.GetInt(`goroutine.timeout`)
	if !VerifyReg(PositiveReg, strconv.Itoa(Goroutine.RoutineTimeOut)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.timeout非法")
	}
	return nil
}

// LoadPortConfig
// @Description: 加载端口配置
// @return       E error
func LoadPortConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadPortConfig-config", util.Strval(r))
		}
	}()
	Port.CenterPort = viper.GetString(`Themis.port`)
	if !VerifyReg(PortReg, Port.CenterPort) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.port端口非法")
	}
	udpTimeOut := viper.GetString(`Themis.UDP-timeout`)
	if !VerifyReg(PositiveReg, udpTimeOut) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.UDP-timeout非法")
	}
	Port.UDPTimeOut, _ = strconv.Atoi(udpTimeOut)
	return nil
}

// LoadServerConfig
// @Description: 加载服务相关配置
// @return       E error
func LoadServerConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadServerConfig-config", util.Strval(r))
		}
	}()
	ServerRegister.ServerModelQueueNum = viper.GetInt(`Themis.server.model-queue`)
	if !VerifyReg(PositiveReg, strconv.Itoa(ServerRegister.ServerModelQueueNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-queue非法")
	}
	ServerRegister.ServerModelHandleNum = viper.GetInt(`Themis.server.model-handle-number`)
	if !VerifyReg(PositiveReg, strconv.Itoa(ServerRegister.ServerModelHandleNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-handle-number")
	}
	ServerBeat.ServerModelBeatEnable = viper.GetBool(`Themis.server.beat-enable`)
	if ServerBeat.ServerModelBeatEnable {
		ServerBeat.ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
		if !VerifyReg(PositiveReg, strconv.Itoa(ServerBeat.ServerModelBeatQueue)) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-queue非法")
		}
		ServerBeat.ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(ServerBeat.ServerBeatTime))) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-time非法")
		}
	} else {
		ServerBeat.ServerModelBeatQueue = 0
		ServerBeat.ServerBeatTime = math.MaxInt
	}
	CreateLeaderAlgorithm = viper.GetString(`Themis.leader-algorithm`)
	ElectionTimeOut = int64(viper.GetInt(`Themis.election-timeout`))
	if !VerifyReg(PositiveReg, strconv.Itoa(int(ElectionTimeOut))) {
		return exception.NewConfigurationError("LoadServerConfig-config", "Themis.election-timeout非法")
	}
	return nil
}

// LoadDatabaseConfig
// @Description: 加载数据库配置
// @return       E error
func LoadDatabaseConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadDatabaseConfig-config", util.Strval(r))
		}
	}()
	Persistence.PersistenceEnable = viper.GetBool(`Themis.database.enable`)
	if Persistence.PersistenceEnable {
		Persistence.PersistenceTime = int64(viper.GetInt(`Themis.database.persistence-time`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Persistence.PersistenceTime))) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config", "Themis.database.persistence-time非法")
		}
		Persistence.SoftDeleteEnable = viper.GetBool(`Themis.database.soft-delete-enable`)
		Database.DatabaseType = viper.GetString(`Themis.database.type`)
		if Database.DatabaseType == "mysql" {
			Database.DatabaseHost = viper.GetString(`Themis.database.mysql.host`)
			if !VerifyReg(IpReg, Database.DatabaseHost) && !VerifyReg(localhostReg, Database.DatabaseHost) {
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
		}
	} else {
		Persistence.PersistenceTime = 0
	}
	return nil
}

// LoadListenConfig
// @Description: 加载监听配置
// @return       E error
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

// LoadClusterConfig
// @Description: 加载集群配置
// @return       E error
func LoadClusterConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadClusterConfig-config", util.Strval(r))
		}
	}()
	Cluster.ClusterEnable = viper.GetBool(`Themis.cluster.enable`)
	if Cluster.ClusterEnable {
		Cluster.TrackEnable = viper.GetBool(`Themis.cluster.track-enable`)
		Cluster.ClusterName = viper.GetString(`Themis.cluster.name`)
		Cluster.IP = viper.GetString(`Themis.cluster.ip`)
		if !VerifyReg(IpReg, Cluster.IP) && !VerifyReg(localhostReg, Cluster.IP) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.ip非法")
		}
		Cluster.Port = viper.GetString(`Themis.cluster.port`)
		if !VerifyReg(PortReg, Cluster.Port) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.port非法")
		}
		Cluster.MaxFollowTimeOut = int64(viper.GetInt(`Themis.cluster.max-follow-timeout`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.MaxFollowTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.max-follow-timeout非法")
		}
		Cluster.MinFollowTimeOut = int64(viper.GetInt(`Themis.cluster.min-follow-timeout`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.MinFollowTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.min-follow-timeout非法")
		}
		Cluster.MaxCandidateTimeOut = int64(viper.GetInt(`Themis.cluster.max-candidate-timeout`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.MaxCandidateTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.max-candidate-timeout非法")
		}
		Cluster.MinCandidateTimeOut = int64(viper.GetInt(`Themis.cluster.min-candidate-timeout`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.MinCandidateTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.min-candidate-timeout非法")
		}
		Cluster.UDPTimeOut = int64(viper.GetInt(`Themis.cluster.udp-timeout`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.UDPTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.udp-timeout非法")
		}
		Cluster.UDPQueueNum = viper.GetInt(`Themis.cluster.udp-queue-num`)
		if !VerifyReg(PositiveReg, strconv.Itoa(Cluster.UDPQueueNum)) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.udp-queue-num非法")
		}
		Cluster.LeaderSnapshotSyncTime = int64(viper.GetInt(`Themis.cluster.leader-snapshot-sync-time`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.LeaderSnapshotSyncTime))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-snapshot-sync-time非法")
		}
		Cluster.LeaderHeartbeatTime = int64(viper.GetInt(`Themis.cluster.leader-heartbeat-time`))
		if !VerifyReg(PositiveReg, strconv.Itoa(int(Cluster.LeaderHeartbeatTime))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-heartbeat-time非法")
		}
		Cluster.LeaderQueueNum = viper.GetInt(`Themis.cluster.leader-queue`)
		if !VerifyReg(PositiveReg, strconv.Itoa(Cluster.LeaderQueueNum)) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-queue非法")
		}
		Cluster.Clusters = make([]map[string]string, 0)
		clusters := viper.Get(`Themis.cluster.clusters`)
		for _, cluster := range clusters.([]interface{}) {
			clusterMap := make(map[string]string)
			clusterMap["ip"] = util.Strval(cluster.(map[string]interface{})["ip"])
			if !VerifyReg(IpReg, clusterMap["ip"]) && !VerifyReg(localhostReg, clusterMap["ip"]) {
				return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.clusters.ip非法")
			}
			clusterMap["port"] = util.Strval(cluster.(map[string]interface{})["port"])
			if !VerifyReg(PortReg, clusterMap["port"]) {
				return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.clusters.port非法")
			}
			Cluster.Clusters = append(Cluster.Clusters, clusterMap)
		}
		Cluster.EnableEncryption = viper.GetBool(`Themis.cluster.enable-encryption`)
		if Cluster.EnableEncryption {
			key := viper.GetString(`Themis.cluster.encryption-key`)
			Cluster.EncryptionKey = []byte(key)
		}
	}
	return nil
}
