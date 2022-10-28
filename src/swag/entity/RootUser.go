package entity

// Root 返回的结果
type Root struct {
	Account  string `json:"account"`
	Password string `json:"password"`
	Token    string `json:"token"`
}

// NewRootModel
//
//	@Description: 创建一个新的User实例
//	@return *Root	返回一个User实例
func NewRootModel() *Root {
	return &Root{}
}
