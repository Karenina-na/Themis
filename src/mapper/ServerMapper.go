package mapper

import (
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

func StorageList(List *util.LinkList[entity.ServerModel], Type int, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("StorageList-mapper", util.Strval(r))
		}
	}()
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("StorageList-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	if List.Length() == 0 {
		util.Loglevel(util.Debug, "StorageList-mapper", "数据为空")
		return true, nil
	}
	for i := 0; i < List.Length(); i++ {
		var mapperModel *entity.ServerMapperMode
		switch Type {
		case NORMAL:
			mapperModel = entity.NewServerMapperMode(List.Get(i), NORMAL)
		case DELETE:
			mapperModel = entity.NewServerMapperMode(List.Get(i), DELETE)
		case LEADER:
			mapperModel = entity.NewServerMapperMode(List.Get(i), LEADER)
		}
		if err := tx.Create(mapperModel).Error; err != nil {
			return false, err
		}
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

func Storage(model *entity.ServerModel, Type int, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Storage-mapper", util.Strval(r))
		}
	}()
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("Storage-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	var mapperModel *entity.ServerMapperMode
	switch Type {
	case NORMAL:
		mapperModel = entity.NewServerMapperMode(*model, NORMAL)
	case DELETE:
		mapperModel = entity.NewServerMapperMode(*model, DELETE)
	case LEADER:
		mapperModel = entity.NewServerMapperMode(*model, LEADER)
	}
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

func SelectAllServers() (S1 []entity.ServerModel, S2 []entity.ServerModel, S3 []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("SelectAllServers-mapper", util.Strval(r))
		}
	}()
	err := DB.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return nil, nil, nil, exception.NewDataBaseError("Storage-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	var modelList []entity.ServerMapperMode
	result := DB.Find(&modelList)
	if err := result.Error; err != nil {
		return nil, nil, nil, err
	}
	var List []entity.ServerModel
	var DeleteList []entity.ServerModel
	var LeaderList []entity.ServerModel
	List = make([]entity.ServerModel, 0, 100)
	DeleteList = make([]entity.ServerModel, 0, 100)
	LeaderList = make([]entity.ServerModel, 0, 100)
	var index int64
	for index = 0; index < result.RowsAffected; index++ {
		switch modelList[index].Type {
		case NORMAL:
			model := entity.NewServerModel()
			model.IP = modelList[index].IP
			model.Name = modelList[index].Name
			model.Port = modelList[index].Port
			model.Colony = modelList[index].Colony
			model.Namespace = modelList[index].Namespace
			model.Time = modelList[index].Time
			List = append(List, *model)
		case DELETE:
			model := entity.NewServerModel()
			model.IP = modelList[index].IP
			model.Name = modelList[index].Name
			model.Port = modelList[index].Port
			model.Colony = modelList[index].Colony
			model.Namespace = modelList[index].Namespace
			model.Time = modelList[index].Time
			DeleteList = append(DeleteList, *model)
		case LEADER:
			model := entity.NewServerModel()
			model.IP = modelList[index].IP
			model.Name = modelList[index].Name
			model.Port = modelList[index].Port
			model.Colony = modelList[index].Colony
			model.Namespace = modelList[index].Namespace
			model.Time = modelList[index].Time
			LeaderList = append(LeaderList, *model)
		}
	}
	util.Loglevel(util.Debug, "SelectAllServers-mapper", "数据查询")
	return List, DeleteList, LeaderList, nil
}

func DeleteServer(model *entity.ServerModel, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-mapper", util.Strval(r))
		}
	}()
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("DeleteServer-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	if err := tx.Delete(&entity.ServerMapperMode{}, "IP = ?", model.IP).Error; err != nil {
		return false, err
	}
	util.Loglevel(util.Debug, "DeleteServer-mapper", "数据删除")
	return true, nil
}

func DeleteAllServer(tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteAllServer-mapper", util.Strval(r))
		}
	}()
	err := tx.AutoMigrate(&entity.ServerMapperMode{})
	if err != nil {
		return false, exception.NewUserError("DeleteAllServer-mapper", "数据库挂载索引结构失败-"+err.Error())
	}
	if err = tx.Where("1 = 1").Delete(&entity.ServerMapperMode{}).Error; err != nil {
		return false, err
	}
	util.Loglevel(util.Debug, "DeleteAllServer-mapper", "数据全部删除")
	return true, nil
}

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
