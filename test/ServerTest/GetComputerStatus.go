package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetComputerStatus(router *gin.Engine) {
	var w *httptest.ResponseRecorder
	fmt.Println("GetComputerStatus:==============================================================")
	w = Base.Get("/api/v1/operator/getStatus", router)
	fmt.Println(w.Body.String())
}
