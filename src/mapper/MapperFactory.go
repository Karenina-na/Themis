package mapper

import (
	"Themis/src/entity"
	"Themis/src/entity/util"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

var DB *gorm.DB

func FactoryInit() {
	var err error
	_, err = os.Stat("./db")
	if os.IsNotExist(err) {
		err := os.MkdirAll("./db", os.ModePerm)
		if err != nil {
			util.Loglevel(util.Error, "DatabaseInit", "创建db文件夹错误"+err.Error())
			os.Exit(0)
		}
	}
	DB, err = gorm.Open(sqlite.Open("./db/Themis.db"), &gorm.Config{})
	if err != nil {
		util.Loglevel(util.Error, "DatabaseInit", "数据库初始化失败-"+err.Error())
		os.Exit(0)
	}
	err = DB.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		util.Loglevel(util.Error, "DatabaseInit", "数据库初始化失败-"+err.Error())
		os.Exit(0)
	}
}
