package exception

import (
	"time"
)

type ControllerPanic struct {
	Name    string
	Time    time.Time
	Message string
}

func NewControllerPanic(name string, Message string) *ControllerPanic {
	return &ControllerPanic{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
