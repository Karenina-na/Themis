package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetAllServersTest(router *gin.Engine) {
	var w *httptest.ResponseRecorder
	fmt.Println("GetAllServersTest:==============================================================")
	w = Base.Get("/api/v1/operator/getInstances", router)
	fmt.Println(w.Body.String())
}
