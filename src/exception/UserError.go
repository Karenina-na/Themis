package exception

import (
	"time"
)

type UserError struct {
	Name    string
	Time    time.Time
	Message string
}

func (e UserError) Error() string {
	return e.Name + "-" + e.Message + e.Time.Format("2006-01-02 15:04:05")
}

func NewUserError(name string, Message string) *UserError {
	return &UserError{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
