package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

// LoadRootConfig
//
//	@Description: 加载管理员信息
//	@return E	error
func LoadRootConfig() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadRootConfig-config", util.Strval(r))
		}
	}()
	config.Root.RootAccount = viper.GetString(`root.account`)
	if config.Root.RootAccount == "" {
		config.Root.RootAccount = "root"
	}
	config.Root.RootPassword = viper.GetString(`root.password`)
	if config.Root.RootPassword == "" {
		config.Root.RootPassword = "root"
	}
	config.Root.TokenEnable = viper.GetBool(`root.token-enable`)
	if config.Root.TokenEnable {
		config.Root.TokenExpireTime = viper.GetInt(`root.token-expire`)
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Root.TokenExpireTime)) {
			return exception.NewConfigurationError("LoadRootConfig-config", "root.token-expire非法")
		}
		config.Root.TokenSignKey = viper.GetString(`root.token-sign`)
		if config.Root.TokenSignKey == "" {
			config.Root.TokenSignKey = "root"
		}
	}
	return nil
}
