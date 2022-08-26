package ServerTest

import (
	"Themis/test/Base"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http/httptest"
)

func GetFollowServer(router *gin.Engine) {
	var w1, w2, w3, w4 *httptest.ResponseRecorder
	w1 = Base.Post("/api/v1/message/getServers", model1, router)
	w2 = Base.Post("/api/v1/message/getServers", model2, router)
	w3 = Base.Post("/api/v1/message/getServers", model3, router)
	w4 = Base.Post("/api/v1/message/getServers", model4, router)
	fmt.Println("GetFollowServer:==============================================================")
	fmt.Println(w1.Body.String())
	fmt.Println(w2.Body.String())
	fmt.Println(w3.Body.String())
	fmt.Println(w4.Body.String())
}
