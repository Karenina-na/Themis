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

// DB 数据库连接
var DB *gorm.DB

// InitMapper
// @Description: 初始化数据库
// @return       E error
func InitMapper() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitMapper-mapper", util.Strval(r))
		}
	}()

	// 选择不同类型的数据库
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

// CloseMapper
// @Description: 关闭数据库
// @return       E error
func CloseMapper() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("CloseMapper-mapper", util.Strval(r))
		}
	}()

	// 关闭数据库
	sqlDB, err := DB.DB()
	err = sqlDB.Close()
	if err != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "数据库连接池关闭错误-"+err.Error())
	}
	return nil
}

// MysqlInit
// @Description: 初始化mysql数据库
// @return       E error
func MysqlInit() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("MysqlInit-mapper", util.Strval(r))
		}
	}()
	var err error
	util.Loglevel(util.Debug, "MysqlInit", "连接数据库mysql")

	// 连接数据库
	dsn := config.Database.DatabaseUser + ":" + config.Database.DatabasePassword + "@tcp(" +
		config.Database.DatabaseHost + ":" + config.Database.DatabasePort + ")/" +
		config.Database.DatabaseName + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "mysql数据库初始化失败-"+err.Error())
	}

	// 设置连接池
	sqlDB, err := DB.DB()
	if err != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "mysql数据库连接池初始化失败-"+err.Error())
	}
	sqlDB.SetMaxOpenConns(config.Database.DatabaseMaxOpenConns)
	sqlDB.SetMaxIdleConns(config.Database.DatabaseMaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Database.DatabaseMaxLifetimeConns) * time.Second)
	return nil
}

// SqlLitInit
// @Description: 初始化sqlite数据库
// @return       E error
func SqlLitInit() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("SqlLitInit-mapper", util.Strval(r))
		}
	}()

	// 判断文件是否存在
	var err error
	_, err = os.Stat("./db")
	if os.IsNotExist(err) {
		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return exception.NewDataBaseError("SqlLitInit-mapper", "sqllit创建db文件夹错误"+err.Error())
		}
	}

	// 连接数据库
	util.Loglevel(util.Debug, "SqlLitInit", "连接数据库sqllit")
	DB, err = gorm.Open(sqlite.Open("./db/Themis.db"), &gorm.Config{})
	if err != nil {
		return exception.NewDataBaseError("SqlLitInit-mapper", "sqllit数据库初始化失败-"+err.Error())
	}
	if err != nil {
		return exception.NewDataBaseError("MysqlInit-mapper", "sqllit数据库连接池初始化失败-"+err.Error())
	}
	return nil
}
