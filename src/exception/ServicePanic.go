package exception

import "time"

type ServicePanic struct {
	Name    string
	Time    time.Time
	Message string
}

func NewServicePanic(name string, Message string) *ServicePanic {
	return &ServicePanic{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
