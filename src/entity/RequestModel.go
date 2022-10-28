package entity

// RequestModel 请求的实体
type RequestModel struct {
	Root Root        `json:"root"`
	Data interface{} `json:"data"`
}

// NewRequestModel
//
//	@Description: 创建一个新的RequestModel实例
//	@return *RequestModel	返回一个RequestModel实例
func NewRequestModel() *RequestModel {
	return &RequestModel{}
}
