package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"time"
)

func CancelDeleteInstance(router *gin.Engine) {
	var w1, w2, w3, w4 *httptest.ResponseRecorder
	w1 = Base.Delete("/api/v1/operator/cancelDeleteInstance", model1, router)
	w2 = Base.Delete("/api/v1/operator/cancelDeleteInstance", model2, router)
	w3 = Base.Delete("/api/v1/operator/cancelDeleteInstance", model3, router)
	w4 = Base.Delete("/api/v1/operator/cancelDeleteInstance", model4, router)
	fmt.Println("CancelDeleteInstance:==============================================================")
	fmt.Println(w1.Body.String())
	fmt.Println(w2.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
	time.Sleep(time.Second)
}
