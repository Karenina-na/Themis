package exception

import (
	"Themis/src/entity/util"
	"os"
)

func HandleException(err interface{}) {
	switch E := err.(type) {
	case *ConfigurationPanic:
		configurationHandle(E)
	case *DataBasePanic:
		dataBaseHandle(E)
	case *ServicePanic:
		serviceHandle(E)
	case *ControllerPanic:
		controllerHandle(E)
	default:
		util.Loglevel(util.Error, "未知错误", util.Strval(err))
	}
}

func configurationHandle(err *ConfigurationPanic) {
	util.Loglevel(util.Info, err.Name, err.Message)
	os.Exit(0)
}

func dataBaseHandle(err *DataBasePanic) {
	util.Loglevel(util.Error, err.Name, err.Message)
	os.Exit(0)
}

func serviceHandle(err *ServicePanic) {
	util.Loglevel(util.Info, err.Name, err.Message)
}

func controllerHandle(err *ControllerPanic) {
	util.Loglevel(util.Info, err.Name, err.Message)
}
