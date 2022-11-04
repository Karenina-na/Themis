package syncBean

// LeaderModel 选举出的Leader模型
type LeaderModel struct {
	// LeaderName	leader名
	LeaderName string
	// LeaderAddress 选举出的leader
	LeaderAddress *SyncAddressModel
	// LeaderServicePort leader服务端口
	LeaderServicePort string
}
