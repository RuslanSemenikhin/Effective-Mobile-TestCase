package http

import (
	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/usecase/control"
	"github.com/gin-gonic/gin"
)

func Start() {
	router := gin.Default()
	initializeRoutes(router)
}

func Stop() {

}

func initializeRoutes(eng *gin.Engine) {
	eng.GET("/service/subscription", func(ctx *gin.Context) { control.ListSubscriptions(ctx) })
}
