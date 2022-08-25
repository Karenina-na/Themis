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
	util.Loglevel(util.Debug, "LoadDatabase", "加载数据库原有文件")
	serverModels, deleteServerModels, leaderServerModels, err := mapper.SelectAllServers()
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
	for i := 0; i < len(deleteServerModels); i++ {
		model := &entity.ServerModel{
			IP:        deleteServerModels[i].IP,
			Port:      deleteServerModels[i].Port,
			Name:      deleteServerModels[i].Name,
			Time:      deleteServerModels[i].Time,
			Colony:    deleteServerModels[i].Colony,
			Namespace: deleteServerModels[i].Namespace,
		}
		_, e1 := RegisterServer(model)
		if e1 != nil {
			return e1
		}
		_, e2 := DeleteServer(model)
		if e2 != nil {
			return e2
		}
	}
	for i := 0; i < len(leaderServerModels); i++ {
		model := entity.ServerModel{
			IP:        leaderServerModels[i].IP,
			Port:      leaderServerModels[i].Port,
			Name:      leaderServerModels[i].Name,
			Time:      leaderServerModels[i].Time,
			Colony:    leaderServerModels[i].Colony,
			Namespace: leaderServerModels[i].Namespace,
		}
		Leaders[model.Namespace] = model
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
	util.Loglevel(util.Debug, "Persistence", "创建持久化数据协程")
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
		}, func(tx *gorm.DB) error {
			list := util.NewLinkList[entity.ServerModel]()
			for _, v := range Leaders {
				list.Append(v)
			}
			b, e := mapper.StorageList(list, mapper.LEADER, tx)
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
