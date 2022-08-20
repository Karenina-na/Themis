package exception

import "time"

type ConfigurationPanic struct {
	Name    string
	Time    time.Time
	Message string
}

func NewConfigurationPanic(name string, Message string) *ConfigurationPanic {
	return &ConfigurationPanic{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
