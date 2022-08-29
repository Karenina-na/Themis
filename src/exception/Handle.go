package exception

import (
	"Themis/src/util"
	"os"
)

func HandleException(err interface{}) {
	switch E := err.(type) {
	case *ConfigurationError:
		configurationExHandle(E)
	case *DataBaseError:
		dataBaseExHandle(E)
	case *SystemError:
		systemExHandle(E)
	case *UserError:
		userExHandle(E)
	default:
		util.Loglevel(util.Error, "未知错误", util.Strval(err))
		os.Exit(0)
	}
}

func configurationExHandle(err *ConfigurationError) {
	util.Loglevel(util.Warn, err.Name, err.Message)
	os.Exit(0)
}

func dataBaseExHandle(err *DataBaseError) {
	util.Loglevel(util.Warn, err.Name, err.Message)
	os.Exit(0)
}

func systemExHandle(err *SystemError) {
	util.Loglevel(util.Warn, err.Name, err.Message)
	os.Exit(0)
}

func userExHandle(err *UserError) {
	util.Loglevel(util.Warn, err.Name, err.Message)
}
