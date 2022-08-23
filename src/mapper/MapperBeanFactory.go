package mapper

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitMapper() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("InitMapper-mapper", util.Strval(r))
		}
	}()
	var err error
	_, err = os.Stat("./db")
	if os.IsNotExist(err) {
		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			return exception.NewDataBaseError("DatabaseInit-mapper", "创建db文件夹错误"+err.Error())
		}
	}
	util.Loglevel(util.Debug, "InitMapper", "连接数据库")
	DB, err = gorm.Open(sqlite.Open("./db/Themis.db"), &gorm.Config{})
	if err != nil {
		return exception.NewDataBaseError("DatabaseInit-mapper", "数据库初始化失败-"+err.Error())
	}
	err = DB.AutoMigrate(&entity.ServerMapperMode{})
	util.Loglevel(util.Debug, "InitMapper", "挂载数据库索引结构")
	if err != nil {
		return exception.NewDataBaseError("DatabaseInit-mapper", "数据库初始化失败-"+err.Error())
	}
	return nil
}
