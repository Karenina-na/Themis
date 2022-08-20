package exception

import "time"

type DataBasePanic struct {
	Name    string
	Time    time.Time
	Message string
}

func NewDataBasePanic(name string, Message string) *DataBasePanic {
	return &DataBasePanic{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
