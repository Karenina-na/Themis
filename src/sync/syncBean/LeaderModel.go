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

// NewLeaderModel
//
//	@Description: 创建LeaderModel
//	@return *LeaderModel	LeaderModel
func NewLeaderModel() *LeaderModel {
	return &LeaderModel{
		LeaderName:        "",
		LeaderAddress:     NewSyncAddressModel(),
		LeaderServicePort: "",
	}
}

// SetLeaderModel
//
//	@Description: 设置LeaderModel
//	@receiver L	LeaderModel
//	@param leaderName	leader名
//	@param leaderIP	leaderIP
//	@param leaderPort	leader端口
//	@param leaderServicePort	leader服务端口
func (L *LeaderModel) SetLeaderModel(leaderName string, leaderIP string,
	leaderPort string, leaderServicePort string) {
	L.LeaderName = leaderName
	L.LeaderAddress.IP = leaderIP
	L.LeaderAddress.Port = leaderPort
	L.LeaderServicePort = leaderServicePort
}
