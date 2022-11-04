package LoadConfigFactory

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"github.com/spf13/viper"
	"strconv"
)

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
	config.Persistence.PersistenceEnable = viper.GetBool(`Themis.database.enable`)
	if config.Persistence.PersistenceEnable {
		config.Persistence.PersistenceTime = int64(viper.GetInt(`Themis.database.persistence-time`))
		if !config.VerifyReg(config.PositiveReg, strconv.Itoa(int(config.Persistence.PersistenceTime))) {
			return exception.NewConfigurationError("LoadDatabaseConfig-config", "Themis.database.persistence-time非法")
		}
		config.Persistence.SoftDeleteEnable = viper.GetBool(`Themis.database.soft-delete-enable`)
		config.Database.DatabaseType = viper.GetString(`Themis.database.type`)
		if config.Database.DatabaseType == "mysql" {
			config.Database.DatabaseHost = viper.GetString(`Themis.database.mysql.host`)
			if !config.VerifyReg(config.IpReg, config.Database.DatabaseHost) && !config.VerifyReg(config.LocalhostReg, config.Database.DatabaseHost) {
				return exception.NewConfigurationError("LoadDatabaseConfig-config",
					"Themis.database.mysql.host非法")
			}
			config.Database.DatabasePort = viper.GetString(`Themis.database.mysql.port`)
			if !config.VerifyReg(config.PortReg, config.Database.DatabasePort) {
				return exception.NewConfigurationError("LoadDatabaseConfig-config",
					"Themis.database.mysql.port非法")
			}
			config.Database.DatabaseName = viper.GetString(`Themis.database.mysql.name`)
			config.Database.DatabaseUser = viper.GetString(`Themis.database.mysql.user`)
			config.Database.DatabasePassword = viper.GetString(`Themis.database.mysql.password`)
			config.Database.DatabaseMaxOpenConns = viper.GetInt(`Themis.database.mysql.max-open-conns`)
			if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Database.DatabaseMaxOpenConns)) {
				return exception.NewConfigurationError("LoadDatabaseConfig-config",
					"Themis.database.mysql.max-open-conns非法")
			}
			config.Database.DatabaseMaxIdleConns = viper.GetInt(`Themis.database.mysql.max-idle-conns`)
			if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Database.DatabaseMaxIdleConns)) {
				return exception.NewConfigurationError("LoadDatabaseConfig-config",
					"Themis.database.mysql.max-idle-conns非法")
			}
			config.Database.DatabaseMaxLifetimeConns = viper.GetInt(`Themis.database.mysql.max-conns-lifetime`)
			if !config.VerifyReg(config.PositiveReg, strconv.Itoa(config.Database.DatabaseMaxLifetimeConns)) {
				return exception.NewConfigurationError("LoadDatabaseConfig-config",
					"Themis.database.mysql.max-conns-lifetime非法")
			}
		}
	} else {
		config.Persistence.PersistenceTime = 0
	}
	return nil
}
