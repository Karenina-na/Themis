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
)

func StorageList(List *util.LinkList[entity.ServerModel], Type int, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("StorageList-mapper", util.Strval(r))
		}
	}()
	if List.Length() == 0 {
		util.Loglevel(util.Debug, "StorageList-mapper", "数据为空")
		return true, nil
	}
	for i := 0; i < List.Length(); i++ {
		var mapperModel *entity.ServerMapperMode
		if Type == NORMAL {
			mapperModel = entity.NewServerMapperMode(List.Get(i), NORMAL)
		} else {
			mapperModel = entity.NewServerMapperMode(List.Get(i), DELETE)
		}
		if err := tx.Create(mapperModel).Error; err != nil {
			return false, err
		}
	}
	var t string
	if Type == NORMAL {
		t = "NORMAL"
	} else {
		t = "DELETE"
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
	var mapperModel *entity.ServerMapperMode
	if Type == NORMAL {
		mapperModel = entity.NewServerMapperMode(*model, NORMAL)
	} else {
		mapperModel = entity.NewServerMapperMode(*model, DELETE)
	}
	if err := tx.Create(mapperModel).Error; err != nil {
		return false, err
	}
	var t string
	if Type == NORMAL {
		t = "NORMAL"
	} else {
		t = "DELETE"
	}
	util.Loglevel(util.Debug, "Storage-mapper", "数据存储"+t)
	return true, nil
}

func SelectAllServers() (S1 []entity.ServerModel, S2 []entity.ServerModel, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("SelectAllServers-mapper", util.Strval(r))
		}
	}()
	var modelList []entity.ServerMapperMode
	result := DB.Find(&modelList)
	if err := result.Error; err != nil {
		return nil, nil, err
	}
	var List []entity.ServerModel
	var DeleteList []entity.ServerModel
	List = make([]entity.ServerModel, 0, 100)
	DeleteList = make([]entity.ServerModel, 0, 100)
	var index int64
	for index = 0; index < result.RowsAffected; index++ {
		if modelList[index].Type == NORMAL {
			model := entity.NewServerModel()
			model.IP = modelList[index].IP
			model.Name = modelList[index].Name
			model.Port = modelList[index].Port
			model.Colony = modelList[index].Colony
			model.Namespace = modelList[index].Namespace
			model.Time = modelList[index].Time
			List = append(List, *model)
		} else {
			model := entity.NewServerModel()
			model.IP = modelList[index].IP
			model.Name = modelList[index].Name
			model.Port = modelList[index].Port
			model.Colony = modelList[index].Colony
			model.Namespace = modelList[index].Namespace
			model.Time = modelList[index].Time
			DeleteList = append(List, *model)
		}
	}
	util.Loglevel(util.Debug, "SelectAllServers-mapper", "数据查询")
	return List, DeleteList, nil
}

func DeleteServer(model *entity.ServerModel, tx *gorm.DB) (B bool, E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteServer-mapper", util.Strval(r))
		}
	}()
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
	if err := tx.Where("1 = 1").Delete(&entity.ServerMapperMode{}).Error; err != nil {
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
	util.Loglevel(util.Debug, "Transaction-mapper", "<<===================================")
	util.Loglevel(util.Debug, "Transaction-mapper", "数据库事务执行"+time.Now().Format("2006-01-02 15:04:05"))
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
		return false, exception.NewDataBaseError("数据库事务中执行失败", util.Strval(err.Error()))
	}
	util.Loglevel(util.Debug, "Transaction-mapper", "数据库事务执行完成-"+time.Now().Format("2006-01-02 15:04:05"))
	util.Loglevel(util.Debug, "Transaction-mapper", "===================================>>")
	return true, nil
}
