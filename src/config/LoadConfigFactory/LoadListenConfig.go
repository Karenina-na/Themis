package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

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
	config.ListenTime = int64(viper.GetInt(`Themis.listen.space-time`))
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.ListenTime))) {
		return exception.NewConfigurationError("LoadListenConfig-config", "Themis.listen.space-time非法")
	}
	return nil
}
