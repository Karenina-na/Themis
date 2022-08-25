package entity

type ServerModel struct {
	IP        string `json:"IP"`
	Port      string `json:"port"`
	Name      string `json:"name"`
	Time      string `json:"time"`
	Colony    string `json:"colony"`
	Namespace string `json:"namespace"`
}

func NewServerModel() *ServerModel {
	return &ServerModel{}
}
