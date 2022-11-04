package LoadConfigFactory

import (
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
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
		if err := LoadRootConfig(); err != nil {
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
