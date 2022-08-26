package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
	"time"
)

func RegisterTest(router *gin.Engine) {
	var w1, w2, w3, w4 *httptest.ResponseRecorder
	w1 = Base.Post("/api/v1/message/register", model1, router)
	w2 = Base.Post("/api/v1/message/register", model2, router)
	w3 = Base.Post("/api/v1/message/register", model3, router)
	w4 = Base.Post("/api/v1/message/register", model4, router)
	fmt.Println("RegisterTest:==============================================================")
	fmt.Println(w1.Body.String())
	fmt.Println(w2.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
	time.Sleep(time.Second)
}
