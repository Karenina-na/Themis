package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

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
	config.Goroutine.MaxRoutineNum = viper.GetInt(`goroutine.max-goroutine`)
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Goroutine.MaxRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.max-goroutine非法")
	}
	config.Goroutine.CoreRoutineNum = viper.GetInt(`goroutine.core-goroutine`)
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Goroutine.CoreRoutineNum)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.core-goroutine非法")
	}
	if config.Goroutine.CoreRoutineNum > config.Goroutine.MaxRoutineNum {
		return exception.NewConfigurationError("LoadRoutineConfig-config",
			"goroutine.core-goroutine大于goroutine.max-goroutine")
	}
	config.Goroutine.RoutineTimeOut = viper.GetInt(`goroutine.timeout`)
	if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Goroutine.RoutineTimeOut)) {
		return exception.NewConfigurationError("LoadRoutineConfig-config", "goroutine.timeout非法")
	}
	return nil
}
