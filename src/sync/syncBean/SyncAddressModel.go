package syncBean

// SyncAddressModel 地址模型
type SyncAddressModel struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
}

// NewSyncAddressModel
//
//	@Description: 创建地址模型
//	@return *SyncAddressModel	地址模型
func NewSyncAddressModel() *SyncAddressModel {
	return &SyncAddressModel{
		IP:   "",
		Port: "",
	}
}
