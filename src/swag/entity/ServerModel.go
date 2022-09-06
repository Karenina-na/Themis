package entity

// ServerModel is a struct that contains the information of a server
type ServerModel struct {
	IP        string `json:"IP" gorm:"column:ip"`
	Port      string `json:"port" gorm:"column:port"`
	Name      string `json:"name" gorm:"column:name"`
	Time      string `json:"time" gorm:"column:time"`
	Colony    string `json:"colony" gorm:"column:colony"`
	Namespace string `json:"namespace" gorm:"column:namespace"`
}
