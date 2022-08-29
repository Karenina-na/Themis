package mapper

import (
	"Themis/src/config"
	"Themis/src/exception"
	"Themis/src/util"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"time"
)

var DB *gorm.DB

// InitMapper 初始化数据库
func InitMapper() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitMapper-mapper", util.Strval(r))
		}
	}()
	switch config.Database.DatabaseType {
	case "sqlite":
		if err := SqlLitInit(); err != nil {
			return err
		}
	case "mysql":
		if err := MysqlInit(); err != nil {
			return err
		}
	default:
		if err := SqlLitInit(); err != nil {
			return err
		}
	}
	return nil
}

func MysqlInit() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("MysqlInit-mapper", util.Strval(r))
		}
	}()
	var err error
	util.Loglevel(util.Debug, "MysqlInit", "连接数据库mysql")
	dsn := config.Database.DatabaseUser + ":" + config.Database.DatabasePassword + "@tcp(" +
		config.Database.DatabaseHost + ":" + config.Database.DatabasePort + ")/" +
		config.Database.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "mysql数据库初始化失败-"+err.Error())
	}
	sqlDB, E := DB.DB()
	if E != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "mysql数据库连接池初始化失败-"+err.Error())
	}
	sqlDB.SetMaxOpenConns(config.Database.DatabaseMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Database.DatabaseMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.DatabaseMaxLifetimeConns) * time.Second)
	return nil
}

func SqlLitInit() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SqlLitInit-mapper", util.Strval(r))
		}
	}()
	var err error
	_, err = os.Stat("./db")
	if os.IsNotExist(err) {
		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return exception.NewDataBaseError("SqlLitInit-mapper", "sqllit创建db文件夹错误"+err.Error())
		}
	}
	util.Loglevel(util.Debug, "SqlLitInit", "连接数据库sqllit")
	DB, err = gorm.Open(sqlite.Open("./db/Themis.db"), &gorm.Config{})
	if err != nil {
		return exception.NewDataBaseError("SqlLitInit-mapper", "sqllit数据库初始化失败-"+err.Error())
	}
	return nil
}
