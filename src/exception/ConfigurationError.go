package exception

import "time"

type ConfigurationError struct {
	Name    string
	Time    time.Time
	Message string
}

func (e ConfigurationError) Error() string {
	return e.Name + "-" + e.Message + e.Time.Format("2006-01-02 15:04:05")
}

func NewConfigurationError(name string, Message string) *ConfigurationError {
	return &ConfigurationError{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
