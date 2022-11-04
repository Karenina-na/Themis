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
