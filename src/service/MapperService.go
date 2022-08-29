package service

import (
	"Themis/src/config"
	"Themis/src/entity"
	"Themis/src/exception"
	"Themis/src/mapper"
	"Themis/src/service/Bean"
	"Themis/src/util"
	"gorm.io/gorm"
	"time"
)

// LoadDatabase 加载数据库
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
		model := serverModels[i].Clone()
		_, e := RegisterServer(model)
		if e != nil {
			return e
		}
	}
	for i := 0; i < len(deleteServerModels); i++ {
		model := deleteServerModels[i].Clone()
		Bean.DeleteInstanceList.Append(*model)
	}
	Bean.Leaders.LeaderModelsListRWLock.Lock()
	for i := 0; i < len(leaderServerModels); i++ {
		model := leaderServerModels[i].Clone()
		if Bean.Leaders.LeaderModelsList[model.Namespace] == nil {
			Bean.Leaders.LeaderModelsList[model.Namespace] = make(map[string]entity.ServerModel)
		}
		Bean.Leaders.LeaderModelsList[model.Namespace][model.Colony] = *model
	}
	Bean.Leaders.LeaderModelsListRWLock.Unlock()
	return nil
}

// Persistence 持久化数据
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
			b, e := mapper.StorageList(Bean.InstanceList, mapper.NORMAL, tx)
			if e != nil || b != true {
				return e
			}
			return nil
		}, func(tx *gorm.DB) error {
			b, e := mapper.StorageList(Bean.DeleteInstanceList, mapper.DELETE, tx)
			if e != nil || b != true {
				return e
			}
			return nil
		}, func(tx *gorm.DB) error {
			list := util.NewLinkList[entity.ServerModel]()
			Bean.Leaders.LeaderModelsListRWLock.RLock()
			for _, v := range Bean.Leaders.LeaderModelsList {
				for _, s := range v {
					list.Append(s)
				}
			}
			b, e := mapper.StorageList(list, mapper.LEADER, tx)
			if e != nil || b != true {
				return e
			}
			Bean.Leaders.LeaderModelsListRWLock.RUnlock()
			return nil
		})
		if e != nil {
			return e
		}
	}
}
