package exception

import "time"

type DataBaseError struct {
	Name    string
	Time    time.Time
	Message string
}

func (e DataBaseError) Error() string {
	return e.Name + "-" + e.Message + e.Time.Format("2006-01-02 15:04:05")
}

func NewDataBaseError(name string, Message string) *DataBaseError {
	return &DataBaseError{
		Name:    name,
		Message: Message,
		Time:    time.Now(),
	}
}
