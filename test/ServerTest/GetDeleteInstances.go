package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetDeleteInstances(router *gin.Engine) {
	var w *httptest.ResponseRecorder
	fmt.Println("GetDeleteInstance:==============================================================")
	w = Base.Get("/api/v1/operator/getDeleteInstance", router)
	fmt.Println(w.Body.String())
}
