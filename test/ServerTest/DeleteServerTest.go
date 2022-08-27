package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"time"
)

func DeleteServerTest(router *gin.Engine) {
	var w1, w2, w3, w4 *httptest.ResponseRecorder
	fmt.Println("DeleteServerTest:==============================================================")
	w1 = Base.Delete("/api/v1/operator/deleteInstance", model1, router)
	w2 = Base.Delete("/api/v1/operator/deleteInstance", model2, router)
	w3 = Base.Delete("/api/v1/operator/deleteInstance", model3, router)
	w4 = Base.Delete("/api/v1/operator/deleteInstance", model4, router)
	fmt.Println(w1.Body.String())
	fmt.Println(w2.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
	time.Sleep(time.Second)
}
