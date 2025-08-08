package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/usecase/control"
	"github.com/gin-gonic/gin"
)

func Start(port int) {
	router := gin.Default()
	initializeRoutes(router)
	portStr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(portStr, router); err != nil {
		log.Fatalf("broken endpoint")
	}
}

func Stop() {

}

func initializeRoutes(eng *gin.Engine) {
	eng.GET("/service/subscription", func(ctx *gin.Context) { control.ListSubscriptions(ctx) })
}
