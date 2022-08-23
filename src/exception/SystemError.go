package exception

import "time"

type SystemError struct {
	Name    string
	Time    time.Time
	Message string
}

func (e SystemError) Error() string {
	return e.Name + "-" + e.Message + e.Time.Format("2006-01-02 15:04:05")
}

func NewSystemError(name string, Message string) *SystemError {
	return &SystemError{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
