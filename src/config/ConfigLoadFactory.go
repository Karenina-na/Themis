package config

import (
	"Themis/src/entity/util"
	"github.com/spf13/viper"
	"os"
	"strconv"
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err == nil {
		MaxRoutineNum = viper.GetInt(`goroutine.MaxRoutineNum`)
		if MaxRoutineNum < 1 {
			util.Loglevel(util.Info, "InitConfig", "MaxRoutineNum非法")
			os.Exit(0)
		}
		CoreRoutineNum = viper.GetInt(`goroutine.CoreRoutineNum`)
		if CoreRoutineNum < 1 {
			util.Loglevel(util.Info, "InitConfig", "CoreRoutineNum非法")
			os.Exit(0)
		}
		if CoreRoutineNum > MaxRoutineNum {
			util.Loglevel(util.Info, "InitConfig", "MaxRoutineNum不能比CoreRoutineNum小")
			os.Exit(0)
		}
		port := viper.GetInt(`Themis.Port`)
		if port < 0 && port > 65535 {
			util.Loglevel(util.Info, "InitConfig", "Port端口错误")
			os.Exit(0)
		}
		Port = strconv.Itoa(port)
		udpPort := viper.GetInt(`Themis.UDPPort`)
		if udpPort < 0 && udpPort > 65535 {
			util.Loglevel(util.Info, "InitConfig", "UDPPort端口错误")
			os.Exit(0)
		} else if udpPort == port {
			util.Loglevel(util.Info, "InitConfig", "Port与UDPPort端口冲突")
			os.Exit(0)
		}
		UDPPort = strconv.Itoa(udpPort)
		ServerModelQueueNum = viper.GetInt(`Themis.server.ServerModelQueueNum`)
		if ServerModelQueueNum <= 0 {
			util.Loglevel(util.Info, "InitConfig", "ServerModelQueueNum非法")
			os.Exit(0)
		}
		ServerModelBeatQueue = viper.GetInt(`Themis.server.ServerModelBeatQueue`)
		if ServerModelBeatQueue <= 0 {
			util.Loglevel(util.Info, "InitConfig", "ServerModelBeatQueue非法")
			os.Exit(0)
		}
		ServerBeatTime = int64(viper.GetInt(`Themis.server.ServerBeatTime`))
		if ServerModelBeatQueue <= 0 {
			util.Loglevel(util.Info, "InitConfig", "ServerBeatTime非法")
			os.Exit(0)
		}
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
