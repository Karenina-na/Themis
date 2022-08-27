package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func ServerBeatTest(router *gin.Engine) {
	var w1, w2, w3, w4 *httptest.ResponseRecorder
	fmt.Println("RegisterTest:==============================================================")
	w1 = Base.Put("/api/v1/message/beat", model1, router)
	w2 = Base.Put("/api/v1/message/beat", model2, router)
	w3 = Base.Put("/api/v1/message/beat", model3, router)
	w4 = Base.Put("/api/v1/message/beat", model4, router)
	fmt.Println(w1.Body.String())
	fmt.Println(w2.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
}
