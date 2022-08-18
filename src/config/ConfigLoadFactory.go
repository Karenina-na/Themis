package config

import (
	"Themis/src/entity/util"
	"github.com/spf13/viper"
	"os"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err == nil {
		MaxRoutineNum = viper.GetInt(`goroutine.MaxRoutineNum`)
		CoreRoutineNum = viper.GetInt(`goroutine.CoreRoutineNum`)
		Port = viper.GetString(`Themis.Port`)
		UDPPort = viper.GetString(`Themis.UDPPort`)
		ServerModelQueueNum = viper.GetInt(`Themis.server.ServerModelQueueNum`)
		ServerModelBeatQueue = viper.GetInt(`Themis.server.ServerModelBeatQueue`)
		ServerBeatTime = int64(viper.GetInt(`Themis.server.ServerBeatTime`))
		CreateLeaderAlgorithm = viper.GetString(`Themis.CreateLeaderAlgorithm`)
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			util.Loglevel(util.Info, "InitConfig", "配置文件config.yaml不存在-"+err.Error())
		} else {
			util.Loglevel(util.Warn, "InitConfig", "未知错误-"+err.Error())
		}
		os.Exit(0)
	}
}
