package entity

import (
	"gorm.io/gorm"
)

// ServerMapperMode 服务器映射表
type ServerMapperMode struct {
	gorm.Model
	ServerModel
	Type int `json:"Type" gorm:"column:type"`
}

// NewServerMapperMode
// @Description: 创建服务器映射表
// @param        model 服务器映射表
// @param        Type  服务器类型
// @return       *ServerMapperMode
func NewServerMapperMode(model ServerModel, Type int) *ServerMapperMode {
	return &ServerMapperMode{
		gorm.Model{},
		ServerModel{
			IP:        model.IP,
			Port:      model.Port,
			Name:      model.Name,
			Time:      model.Time,
			Colony:    model.Colony,
			Namespace: model.Namespace,
		},
		Type,
	}
}

// TableName
// @Description: 获取表名
// @receiver     model  服务器映射表
// @return       string 表名
func (model ServerMapperMode) TableName() string {
	return "tb_server_mapper"
}

// UnPack
//
//	@Description: 解包
//	@receiver model	服务器映射表
//	@return *entity.ServerModel	服务器模型
func (model ServerMapperMode) UnPack() *ServerModel {
	return &ServerModel{
		IP:        model.IP,
		Port:      model.Port,
		Name:      model.Name,
		Time:      model.Time,
		Colony:    model.Colony,
		Namespace: model.Namespace,
	}
}
