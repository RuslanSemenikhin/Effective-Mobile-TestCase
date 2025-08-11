package control

import (
	"log"
	"net/http"

	g "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AddSubRequest struct {
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	UserUuid    string `json:"user_uuid"`
	StartDate   string `json:"start_date"`
	StopDate    string `json:"stop_date"`
}

func ListSubscriptions(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'ListSubscriptions' into endpoint service start with requestID - '%s'", reqID)
	dateAt := ctx.Query("start_date")
	dateTo := ctx.Query("stop_date")
	userUuid := ctx.Query("user_uuid")
	serviceName := ctx.Query("service_name")
	req := &g.ListSubscriptionsRequest{
		ReqId:       reqID,
		StartDate:   &dateAt,
		StopDate:    &dateTo,
		UserUuid:    &userUuid,
		ServiceName: &serviceName,
	}
	subscriptions, err := client.ListSubscriptions(ctx, req)
	if err != nil {
		log.Printf("function 'ListSubscriptions' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
	log.Printf("function 'ListSubscriptions' into endpoint service finished successfuly with requestID - '%s'", reqID)
}

func AddSubscription(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'AddSubscription' into endpoint service start with requestID - '%s'", reqID)

	var req AddSubRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("function 'AddSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	s := &g.Subscription{
		ServiceName: req.ServiceName,
		Price:       int64(req.Price),
		UserUuid:    req.UserUuid,
		StartDate:   req.StartDate,
		StopDate:    req.StopDate,
	}

	r := &g.AddSubscriptionRequest{
		ReqID:        reqID,
		Subscription: s,
	}

	resp, err := client.AddSubscription(ctx, r)
	if err != nil {
		log.Printf("function 'AddSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, resp)
	log.Printf("function 'AddSubscription' into endpoint service finished successfuly with requestID - '%s'", reqID)
}
