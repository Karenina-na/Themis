package Init

import (
	"Themis/src/exception"
	"Themis/src/service"
)

func ThemisCloseFactory() {
	err := service.Close()
	if err != nil {
		exception.HandleException(exception.NewSystemError("ThemisCloseFactory", err.Error()))
	}
}
