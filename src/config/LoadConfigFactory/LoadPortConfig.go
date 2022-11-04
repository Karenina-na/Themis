package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

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
	config.Port.CenterPort = viper.GetString(`Themis.port`)
	if !config.VerifyReg(config.PortReg, config.Port.CenterPort) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.port端口非法")
	}
	udpTimeOut := viper.GetString(`Themis.UDP-timeout`)
	if !config.VerifyReg(config.PositiveReg, udpTimeOut) {
		return exception.NewConfigurationError("LoadPortConfig-config", "Themis.UDP-timeout非法")
	}
	config.Port.UDPTimeOut, _ = strconv.Atoi(udpTimeOut)
	return nil
}
