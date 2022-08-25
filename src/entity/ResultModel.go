package entity

const (
	Success = 20010 + iota
	False
)

type ResultModel struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func NewSuccessResult(data interface{}) *ResultModel {
	return &ResultModel{
		Code: Success,
		Data: data,
	}
}

func NewFalseResult(data interface{}, Message string) *ResultModel {
	return &ResultModel{
		Code:    False,
		Data:    data,
		Message: Message,
	}
}
