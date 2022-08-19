package entity

import "gorm.io/gorm"

type ServerMapperMode struct {
	gorm.Model
	ServerModel
	Type int `json:"Type"`
}

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