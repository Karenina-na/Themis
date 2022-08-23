package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/util"
	"gorm.io/gorm"
	"time"
)

func LoadDatabase() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadDatabase-service", util.Strval(r))
		}
	}()
	serverModels, deleteServerModels, err := mapper.SelectAllServers()
	if err != nil {
		return err
	}
	for i := 0; i < len(serverModels); i++ {
		model := &entity.ServerModel{
			IP:        serverModels[i].IP,
			Port:      serverModels[i].Port,
			Name:      serverModels[i].Name,
			Time:      serverModels[i].Time,
			Colony:    serverModels[i].Colony,
			Namespace: serverModels[i].Namespace,
		}
		_, e := RegisterServer(model)
		if e != nil {
			return e
		}
	}
	for i := 0; i < len(deleteServerModels)-1; i++ {
		model := &entity.ServerModel{
			IP:        serverModels[i].IP,
			Port:      serverModels[i].Port,
			Name:      serverModels[i].Name,
			Time:      serverModels[i].Time,
			Colony:    serverModels[i].Colony,
			Namespace: serverModels[i].Namespace,
		}
		_, e := DeleteServer(model)
		if e != nil {
			return e
		}
	}
	return nil
}

func Persistence() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Persistence-service", util.Strval(r))
		}
	}()
	for {
		time.Sleep(time.Duration(config.PersistenceTime) * time.Second)
		_, e := mapper.Transaction(func(tx *gorm.DB) error {
			b, e := mapper.DeleteAllServer(tx)
			if e != nil || b != true {
				return e
			}
			return nil
		}, func(tx *gorm.DB) error {
			b, e := mapper.StorageList(InstanceList, mapper.NORMAL, tx)
			if e != nil || b != true {
				return e
			}
			return nil
		}, func(tx *gorm.DB) error {
			b, e := mapper.StorageList(DeleteInstanceList, mapper.DELETE, tx)
			if e != nil || b != true {
				return e
			}
			return nil
		})
		if e != nil {
			return e
		}
	}
}

func DeleteMapper(model *entity.ServerModel) (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("DeleteMapper-service", util.Strval(r))
		}
	}()
	_, e := mapper.Transaction(func(tx *gorm.DB) error {
		b, e := mapper.DeleteServer(model, tx)
		if e != nil || b != true {
			return e
		}
		return nil
	})
	if e != nil {
		return e
	}
	return nil
}
