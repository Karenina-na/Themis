package ServerTest

import "Themis/src/entity"

var model1 = entity.ServerModel{
	IP:        "111.111.111.111",
	Port:      "80",
	Name:      "Server-A",
	Time:      "2022-8-8",
	Colony:    "Default",
	Namespace: "Default",
}
var model2 = entity.ServerModel{
	IP:        "222.222.222.222",
	Port:      "80",
	Name:      "Server-B",
	Time:      "2022-8-8",
	Colony:    "Default",
	Namespace: "Default",
}
var model3 = entity.ServerModel{
	IP:        "333.333.333.333",
	Port:      "80",
	Name:      "Server-C",
	Time:      "2022-8-8",
	Colony:    "A",
	Namespace: "A",
}
var model4 = entity.ServerModel{
	IP:        "444.444.444.444",
	Port:      "80",
	Name:      "Server-D",
	Time:      "2022-8-8",
	Colony:    "A",
	Namespace: "B",
}
