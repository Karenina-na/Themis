package mapper

import (
	"Themis/src/entity"
	"Themis/src/exception"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func InitMapper() (E any) {
	defer func() {
		E = recover()
	}()
	var err error
	_, err = os.Stat("./db")
	if os.IsNotExist(err) {
		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			panic(exception.NewDataBasePanic("DatabaseInit", "创建db文件夹错误"+err.Error()))
		}
	}
	DB, err = gorm.Open(sqlite.Open("./db/Themis.db"), &gorm.Config{})
	if err != nil {
		panic(exception.NewDataBasePanic("DatabaseInit", "数据库初始化失败-"+err.Error()))
	}
	err = DB.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		panic(exception.NewDataBasePanic("DatabaseInit", "数据库初始化失败-"+err.Error()))
	}
	return nil
}
