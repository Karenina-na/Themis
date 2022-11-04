package mapper

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/util"
	"gorm.io/gorm"
	"time"
)

const (
	NORMAL = iota
	DELETE
	LEADER
)

// StorageList
// @Description: 存储服务列表
// @param        List 服务列表
// @param        Type 存储类型
// @param        tx   事务
// @return       B    是否成功
// @return       E    错误信息
func StorageList(List *util.LinkList[entity.ServerModel], Type int, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("StorageList-mapper", util.Strval(r))
		}
	}()

	// 绑定orm
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("StorageList-mapper", "数据库挂载索引结构失败-"+err.Error())
	}

	// 判空，减少io
	if List.Length() == 0 {
		util.Loglevel(util.Debug, "StorageList-mapper", "数据为空")
		return true, nil
	}

	// 遍历存储
	var ERR error
	List.Iterator(func(index int, model entity.ServerModel) {

		// 按类型包装数据
		var mapperModel *entity.ServerMapperMode
		switch Type {
		case NORMAL:
			mapperModel = entity.NewServerMapperMode(model, NORMAL)
		case DELETE:
			mapperModel = entity.NewServerMapperMode(model, DELETE)
		case LEADER:
			mapperModel = entity.NewServerMapperMode(model, LEADER)
		}

		// 存储
		if err := tx.Create(mapperModel).Error; err != nil {
			ERR = err
			return
		}
	})
	if ERR != nil {
		return false, ERR
	}
	var t string
	switch Type {
	case NORMAL:
		t = "NORMAL"
	case DELETE:
		t = "DELETE"
	case LEADER:
		t = "LEADER"
	}
	util.Loglevel(util.Debug, "StorageList-mapper", "数据批量存储-"+t)
	return true, nil
}

// Storage
// @Description: 存储服务
// @param        model 服务
// @param        Type  存储类型
// @param        tx    事务
// @return       B     是否成功
// @return       E     错误信息
func Storage(model *entity.ServerModel, Type int, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Storage-mapper", util.Strval(r))
		}
	}()

	// 绑定orm
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("Storage-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	var mapperModel *entity.ServerMapperMode

	// 按类型包装数据
	switch Type {
	case NORMAL:
		mapperModel = entity.NewServerMapperMode(*model, NORMAL)
	case DELETE:
		mapperModel = entity.NewServerMapperMode(*model, DELETE)
	case LEADER:
		mapperModel = entity.NewServerMapperMode(*model, LEADER)
	}

	// 存储
	if err := tx.Create(mapperModel).Error; err != nil {
		return false, err
	}
	var t string
	switch Type {
	case NORMAL:
		t = "NORMAL"
	case DELETE:
		t = "DELETE"
	case LEADER:
		t = "LEADER"
	}
	util.Loglevel(util.Debug, "Storage-mapper", "数据存储"+t)
	return true, nil
}

// SelectAllServers
// @Description: 查询所有服务
// @return       S1 所有服务列表
// @return       S2 黑名单服务列表
// @return       S3 领导者服务列表
// @return       E  错误信息
func SelectAllServers() (S1 []entity.ServerModel, S2 []entity.ServerModel, S3 []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("SelectAllServers-mapper", util.Strval(r))
		}
	}()

	// 绑定orm
	err := DB.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return nil, nil, nil, exception.NewDataBaseError("Storage-mapper", "数据库挂载索引结构失败-"+err.Error())
	}

	// 查询
	var modelList []entity.ServerMapperMode
	result := DB.Find(&modelList)
	if err := result.Error; err != nil {
		return nil, nil, nil, err
	}

	// 创建解包存储数据
	var List []entity.ServerModel
	var DeleteList []entity.ServerModel
	var LeaderList []entity.ServerModel
	List = make([]entity.ServerModel, 0, 100)
	DeleteList = make([]entity.ServerModel, 0, 100)
	LeaderList = make([]entity.ServerModel, 0, 100)
	var index int64

	// 遍历数据
	for index = 0; index < result.RowsAffected; index++ {
		model := *modelList[index].UnPack()
		switch modelList[index].Type {
		case NORMAL:
			List = append(List, model)
		case DELETE:
			DeleteList = append(DeleteList, model)
		case LEADER:
			LeaderList = append(LeaderList, model)
		}
	}
	util.Loglevel(util.Debug, "SelectAllServers-mapper", "数据查询")
	return List, DeleteList, LeaderList, nil
}

// DeleteServer
// @Description: 删除服务
// @param        model 服务
// @param        tx    事务
// @return       B     是否成功
// @return       E     错误信息
func DeleteServer(model *entity.ServerModel, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-mapper", util.Strval(r))
		}
	}()

	// 绑定orm
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("DeleteServer-mapper", "数据库挂载索引结构失败-"+err.Error())
	}

	// 判断软删除
	if config.Persistence.SoftDeleteEnable {
		if err := tx.Delete(&entity.ServerMapperMode{}, "IP = ?", model.IP).Error; err != nil {
			return false, err
		}
	} else {
		if err := tx.Unscoped().Delete(&entity.ServerMapperMode{}, "IP = ?", model.IP).Error; err != nil {
			return false, err
		}
	}
	util.Loglevel(util.Debug, "DeleteServer-mapper", "数据删除")
	return true, nil
}

// DeleteAllServer
// @Description: 删除所有服务
// @param        tx 事务
// @return       B  是否成功
// @return       E  错误信息
func DeleteAllServer(tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteAllServer-mapper", util.Strval(r))
		}
	}()

	// 绑定orm
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("DeleteAllServer-mapper", "数据库挂载索引结构失败-"+err.Error())
	}

	// 判断软删除
	if config.Persistence.SoftDeleteEnable {
		if err = tx.Where("1 = 1").Delete(&entity.ServerMapperMode{}).Error; err != nil {
			return false, err
		}
	} else {
		if err = tx.Unscoped().Where("1 = 1").Delete(&entity.ServerMapperMode{}).Error; err != nil {
			return false, err
		}
	}
	util.Loglevel(util.Debug, "DeleteAllServer-mapper", "数据全部删除")
	return true, nil
}

// Transaction
// @Description: 事务
// @param        List 服务列表
// @return       B    是否成功
// @return       E    错误信息
func Transaction(List ...func(tx *gorm.DB) error) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Transaction-mapper", util.Strval(r))
		}
	}()
	t := time.Now().Format("2006-01-02 15:04:05")
	util.Loglevel(util.Debug, "Transaction-mapper", "<<===================================")
	util.Loglevel(util.Debug, "Transaction-mapper", "数据库事务执行"+t)

	// 事务组
	err := DB.Transaction(func(tx *gorm.DB) error {
		for _, f := range List {
			err := f(tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return false, exception.NewDataBaseError("Transaction-mapper", "数据库事务中执行失败 "+util.Strval(err.Error())+" "+t)
	}
	util.Loglevel(util.Debug, "Transaction-mapper", "数据库事务执行完成-"+time.Now().Format("2006-01-02 15:04:05"))
	util.Loglevel(util.Debug, "Transaction-mapper", "===================================>>")
	return true, nil
}
