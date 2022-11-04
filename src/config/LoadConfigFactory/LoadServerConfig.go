package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"math"
	"strconv"
)

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
	config.ServerRegister.ServerModelQueueNum = viper.GetInt(`Themis.server.model-queue`)
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.ServerRegister.ServerModelQueueNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-queue非法")
	}
	config.ServerRegister.ServerModelHandleNum = viper.GetInt(`Themis.server.model-handle-number`)
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.ServerRegister.ServerModelHandleNum)) {
		return exception.NewConfigurationError("LoadServerConfig-config",
			"Themis.server.model-handle-number")
	}
	config.ServerBeat.ServerModelBeatEnable = viper.GetBool(`Themis.server.beat-enable`)
	if config.ServerBeat.ServerModelBeatEnable {
		config.ServerBeat.ServerModelBeatQueue = viper.GetInt(`Themis.server.beat-queue`)
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.ServerBeat.ServerModelBeatQueue)) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-queue非法")
		}
		config.ServerBeat.ServerBeatTime = int64(viper.GetInt(`Themis.server.beat-time`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.ServerBeat.ServerBeatTime))) {
			return exception.NewConfigurationError("LoadServerConfig-config", "Themis.server.beat-time非法")
		}
	} else {
		config.ServerBeat.ServerModelBeatQueue = 0
		config.ServerBeat.ServerBeatTime = math.MaxInt
	}
	config.CreateLeaderAlgorithm = viper.GetString(`Themis.leader-algorithm`)
	config.ElectionTimeOut = int64(viper.GetInt(`Themis.election-timeout`))
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.ElectionTimeOut))) {
		return exception.NewConfigurationError("LoadServerConfig-config", "Themis.election-timeout非法")
	}
	return nil
}
