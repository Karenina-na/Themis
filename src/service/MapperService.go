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

// LoadDatabase
// @Description: 加载数据库数据
// @return       E error
func LoadDatabase() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewSystemError("LoadDatabase-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "LoadDatabase", "加载数据库原有文件")

	//加载数据库文件
	serverModels, deleteServerModels, leaderServerModels, err := mapper.SelectAllServers()
	if err != nil {
		return err
	}

	//注册服务
	for i := 0; i < len(serverModels); i++ {
		model := serverModels[i].Clone()
		B, e := RegisterServer(model)
		if e != nil || B != true {
			return e
		}
	}

	//注册黑名单服务
	for i := 0; i < len(deleteServerModels); i++ {
		model := deleteServerModels[i].Clone()
		Bean.DeleteInstanceList.Append(*model)
	}

	//注册领导者服务
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

// Persistence
// @Description: 持久化数据
// @return       E error
func Persistence() (E error) {
	defer func() {
		r := recover()
		if r != nil {
			E = exception.NewUserError("Persistence-service", util.Strval(r))
		}
	}()
	util.Loglevel(util.Debug, "Persistence", "创建持久化数据协程")
	for {
		select {
		case <-Bean.ServiceCloseChan:
			util.Loglevel(util.Debug, "Persistence", "持久化数据协程退出")
			return nil
		case <-time.After(time.Second * time.Duration(config.Persistence.PersistenceTime)):

			//事务
			B, e := mapper.Transaction(func(tx *gorm.DB) error {

				//删除所有已经持久化的数据信息
				b, e := mapper.DeleteAllServer(tx)
				if e != nil || b != true {
					return e
				}
				return nil
			}, func(tx *gorm.DB) error {

				//持久化实例信息
				b, e := mapper.StorageList(Bean.InstanceList, mapper.NORMAL, tx)
				if e != nil || b != true {
					return e
				}
				return nil
			}, func(tx *gorm.DB) error {

				//持久化黑名单信息
				b, e := mapper.StorageList(Bean.DeleteInstanceList, mapper.DELETE, tx)
				if e != nil || b != true {
					return e
				}
				return nil
			}, func(tx *gorm.DB) error {

				//持久化领导者信息
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
			if e != nil || B != true {
				return e
			}
		}
	}
}
