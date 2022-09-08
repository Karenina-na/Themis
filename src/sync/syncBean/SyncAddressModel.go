package syncBean

// SyncAddressModel 地址模型
type SyncAddressModel struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

func NewSyncAddressModel() *SyncAddressModel {
	return &SyncAddressModel{
		IP:   "",
		Port: "",
	}
}
