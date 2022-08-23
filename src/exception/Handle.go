package exception

import (
	util2 "Themis/src/util"
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
		util2.Loglevel(util2.Error, "未知错误", util2.Strval(err))
		os.Exit(0)
	}
}

func configurationExHandle(err *ConfigurationError) {
	util2.Loglevel(util2.Error, err.Name, err.Message+"-"+err.Error()+"-"+err.Time.Format("2006-01-02 15:04:05"))
	os.Exit(0)
}

func dataBaseExHandle(err *DataBaseError) {
	util2.Loglevel(util2.Error, err.Name, err.Message+"-"+err.Error()+"-"+err.Time.Format("2006-01-02 15:04:05"))
	os.Exit(0)
}

func systemExHandle(err *SystemError) {
	util2.Loglevel(util2.Error, err.Name, err.Message+"-"+err.Error()+"-"+err.Time.Format("2006-01-02 15:04:05"))
	os.Exit(0)
}

func userExHandle(err *UserError) {
	util2.Loglevel(util2.Warn, err.Name, err.Message+"-"+err.Error()+"-"+err.Time.Format("2006-01-02 15:04:05"))
}
