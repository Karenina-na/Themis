package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetInstancesByQuery(router *gin.Engine) {
	var w *httptest.ResponseRecorder
	fmt.Println("GetInstancesByQuery:==============================================================")
	w = Base.Post("/api/v1/operator/getInstances", model1, router)
	fmt.Println(w.Body.String())
}
