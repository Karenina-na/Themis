package entity

type ServerModel struct {
	IP        string `json:"IP"`
	Port      string `json:"Port"`
	Name      string `json:"Name"`
	Time      string `json:"Time"`
	Colony    string `json:"Colony"`
	Namespace string `json:"Namespace"`
}

func NewServerModel() *ServerModel {
	return &ServerModel{}
}
