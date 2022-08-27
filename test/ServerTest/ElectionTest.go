package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"time"
)

func ElectionTest(router *gin.Engine) {
	var w1, w3, w4 *httptest.ResponseRecorder
	fmt.Println("ElectionTest:==============================================================")
	w1 = Base.Put("/api/v1/message/election", model1, router)
	w3 = Base.Put("/api/v1/message/election", model3, router)
	w4 = Base.Put("/api/v1/message/election", model4, router)
	fmt.Println(w1.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
	time.Sleep(time.Second)
}
