package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

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
	config.Cluster.ClusterEnable = viper.GetBool(`Themis.cluster.enable`)
	if config.Cluster.ClusterEnable {
		config.Cluster.TrackEnable = viper.GetBool(`Themis.cluster.track-enable`)
		config.Cluster.ClusterName = viper.GetString(`Themis.cluster.name`)
		config.Cluster.IP = viper.GetString(`Themis.cluster.ip`)
		if !config.VerifyReg(config.IpReg, config.Cluster.IP) && !config.VerifyReg(config.LocalhostReg, config.Cluster.IP) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.ip非法")
		}
		config.Cluster.Port = viper.GetString(`Themis.cluster.port`)
		if !config.VerifyReg(config.PortReg, config.Cluster.Port) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.port非法")
		}
		config.Cluster.MaxFollowTimeOut = int64(viper.GetInt(`Themis.cluster.max-follow-timeout`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.MaxFollowTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.max-follow-timeout非法")
		}
		config.Cluster.MinFollowTimeOut = int64(viper.GetInt(`Themis.cluster.min-follow-timeout`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.MinFollowTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.min-follow-timeout非法")
		}
		config.Cluster.MaxCandidateTimeOut = int64(viper.GetInt(`Themis.cluster.max-candidate-timeout`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.MaxCandidateTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.max-candidate-timeout非法")
		}
		config.Cluster.MinCandidateTimeOut = int64(viper.GetInt(`Themis.cluster.min-candidate-timeout`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.MinCandidateTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.min-candidate-timeout非法")
		}
		config.Cluster.UDPTimeOut = int64(viper.GetInt(`Themis.cluster.udp-timeout`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.UDPTimeOut))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.udp-timeout非法")
		}
		config.Cluster.UDPQueueNum = viper.GetInt(`Themis.cluster.udp-queue-num`)
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Cluster.UDPQueueNum)) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.udp-queue-num非法")
		}
		config.Cluster.LeaderSnapshotSyncTime = int64(viper.GetInt(`Themis.cluster.leader-snapshot-sync-time`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.LeaderSnapshotSyncTime))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-snapshot-sync-time非法")
		}
		config.Cluster.LeaderHeartbeatTime = int64(viper.GetInt(`Themis.cluster.leader-heartbeat-time`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Cluster.LeaderHeartbeatTime))) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-heartbeat-time非法")
		}
		config.Cluster.LeaderQueueNum = viper.GetInt(`Themis.cluster.leader-queue`)
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Cluster.LeaderQueueNum)) {
			return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.leader-queue非法")
		}
		config.Cluster.Clusters = make([]map[string]string, 0)
		clusters := viper.Get(`Themis.cluster.clusters`)
		for _, cluster := range clusters.([]interface{}) {
			clusterMap := make(map[string]string)
			clusterMap["ip"] = util.Strval(cluster.(map[string]interface{})["ip"])
			if !config.VerifyReg(config.IpReg, clusterMap["ip"]) && !config.VerifyReg(config.LocalhostReg, clusterMap["ip"]) {
				return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.clusters.ip非法")
			}
			clusterMap["port"] = util.Strval(cluster.(map[string]interface{})["port"])
			if !config.VerifyReg(config.PortReg, clusterMap["port"]) {
				return exception.NewConfigurationError("LoadClusterConfig-config", "Themis.cluster.clusters.port非法")
			}
			config.Cluster.Clusters = append(config.Cluster.Clusters, clusterMap)
		}
		config.Cluster.EnableEncryption = viper.GetBool(`Themis.cluster.enable-encryption`)
		if config.Cluster.EnableEncryption {
			key := viper.GetString(`Themis.cluster.encryption-key`)
			config.Cluster.EncryptionKey = []byte(key)
		}
	}
	return nil
}
