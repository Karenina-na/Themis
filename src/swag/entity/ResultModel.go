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

// NewSuccessResult
// @Description: 返回成功的结果
// @param        data 返回的数据
// @return       *ResultModel 返回的结果
func NewSuccessResult(data interface{}) *ResultModel {
	return &ResultModel{
		Code: Success,
		Data: data,
	}
}

// NewFalseResult
// @Description: 返回失败的结果
// @param        data    返回的数据
// @param        SyncMessage 返回的消息
// @return       *ResultModel 返回的结果
func NewFalseResult(data interface{}, Message string) *ResultModel {
	return &ResultModel{
		Code:    False,
		Data:    data,
		Message: Message,
	}
}
