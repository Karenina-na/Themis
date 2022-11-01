package entity

// RequestModel 请求的实体
type RequestModel struct {
	Root Root        `json:"root"`
	Data interface{} `json:"data"`
}
