package mapper

import (
	"Themis/src/entity"
	"Themis/src/entity/util"
)

const (
	NORMAL = iota
	DELETE
)

func StorageList(List *util.LinkList[entity.ServerModel], Type int) bool {
	for i := 0; i < List.Length(); i++ {
		var mapperModel *entity.ServerMapperMode
		if Type == NORMAL {
			mapperModel = entity.NewServerMapperMode(List.Get(i), NORMAL)
		} else {
			mapperModel = entity.NewServerMapperMode(List.Get(i), DELETE)
		}
		DB.Create(mapperModel)
	}
	return true
}

func Storage(model *entity.ServerModel, Type int) bool {
	var mapperModel *entity.ServerMapperMode
	if Type == NORMAL {
		mapperModel = entity.NewServerMapperMode(*model, NORMAL)
	} else {
		mapperModel = entity.NewServerMapperMode(*model, DELETE)
	}
	DB.Create(mapperModel)
	return true
}

func SelectAllServers() ([]entity.ServerModel, []entity.ServerModel) {
	var modelList []entity.ServerMapperMode
	result := DB.Find(&modelList)
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
	return List, DeleteList
}

func DeleteServer(model *entity.ServerModel) bool {
	DB.Delete(&entity.ServerMapperMode{}, "IP = ?", model.IP)
	return true
}

func DeleteAllServer() bool {
	DB.Where("1 = 1").Delete(&entity.ServerMapperMode{})
	return true
}
