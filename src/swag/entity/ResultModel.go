package entity

const (
	Success = 20010 + iota
	False
)

// ResultModel 返回的结果
type ResultModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
