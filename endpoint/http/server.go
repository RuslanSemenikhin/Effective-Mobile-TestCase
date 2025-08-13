package http

import (
	"fmt"
	"log"
	"net/http"

	g "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/usecase/control"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type EndpointServer struct {
	Router     *gin.Engine
	GrpcClient g.SubscriptionServiceClient
}

func NewEndpointServer() *EndpointServer {
	return &EndpointServer{}
}

func (e *EndpointServer) Start(
	endpointPort int,
	controllerHost string,
	controllerPort int,
) {
	address := fmt.Sprintf("%s:%d", controllerHost, controllerPort)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to controller: %v", err)
	}
	e.GrpcClient = g.NewSubscriptionServiceClient(conn)

	router := gin.Default()
	initializeRoutes(router, e.GrpcClient)
	e.Router = router

	e.startHttp(endpointPort)
}

func (e *EndpointServer) startHttp(port int) {
	portStr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(portStr, e.Router); err != nil {
		log.Fatalf("can`t start http-endpoint - '%v'", err.Error())
	}
}

func initializeRoutes(router *gin.Engine, client g.SubscriptionServiceClient) {
	router.GET("/service/subscription", func(ctx *gin.Context) { control.ListSubscriptions(ctx, client) })
	router.POST("/service/subscription", func(ctx *gin.Context) { control.AddSubscription(ctx, client) })
	router.PUT("/service/:service_name/user/:user_uuid/subscription", func(ctx *gin.Context) { control.UpdatedSubscription(ctx, client) })
	router.DELETE("/service/:service_name/user/:user_uuid/subscription", func(ctx *gin.Context) { control.DeleteSubscription(ctx, client) })
	router.GET("/service/subscription/totalprice/start/:date_start/stop/:date_stop", func(ctx *gin.Context) { control.TotalPrice(ctx, client) })
}
