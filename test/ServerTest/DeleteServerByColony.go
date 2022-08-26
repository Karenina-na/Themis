package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func DeleteServerByColony(router *gin.Engine) {
	var w1, w3, w4 *httptest.ResponseRecorder
	w1 = Base.Delete("/api/v1/operator/deleteColony", model1, router)
	w3 = Base.Delete("/api/v1/operator/deleteColony", model3, router)
	w4 = Base.Delete("/api/v1/operator/deleteColony", model4, router)
	fmt.Println("DeleteServerByColony:==========================================================")
	fmt.Println(w1.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
}
