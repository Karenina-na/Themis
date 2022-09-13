package exception

import (
	"time"
)

// UserError 用户错误
type UserError struct {
	Name    string
	Time    time.Time
	Message string
}

// Error
// @Description: 实现error接口
// @receiver     e      UserError
// @return       string 错误信息
func (e UserError) Error() string {
	return e.Name + "-" + e.Message + e.Time.Format("2006-01-02 15:04:05")
}

// NewUserError
// @Description: 创建用户错误
// @param        name    string 错误名称
// @param        SyncMessage string 错误信息
// @return       *UserError 用户错误
func NewUserError(name string, Message string) *UserError {
	return &UserError{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
