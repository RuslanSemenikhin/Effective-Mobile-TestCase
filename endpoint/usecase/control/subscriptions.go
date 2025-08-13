package control

import (
	"log"
	"net/http"

	g "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SubRequest struct {
	ServiceName string `json:"service_name"`
	Price       int64  `json:"price"`
	UserUuid    string `json:"user_uuid"`
	StartDate   string `json:"start_date"`
	StopDate    string `json:"stop_date"`
}

type PropsForUpdateREquest struct {
	Price     *int64  `json:"price"`
	StartDate *string `json:"start_date"`
	StopDate  *string `json:"stop_date"`
}

func ListSubscriptions(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'ListSubscriptions' into endpoint service start with requestID - '%s'", reqID)
	userUuid := ctx.Query("user_uuid")
	serviceName := ctx.Query("service_name")
	req := &g.ListSubscriptionsRequest{
		ReqId:       reqID,
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

func TotalPrice(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'TotalPrice' into endpoint service start with requestID - '%s'", reqID)
	startDate := ctx.Param("date_start")
	stopDate := ctx.Param("date_stop")
	userUuid := ctx.Query("user_uuid")
	serviceName := ctx.Query("service_name")

	req := &g.TotalPriceRequest{
		ReqID:       reqID,
		StartDate:   startDate,
		StopDate:    stopDate,
		UserUuid:    &userUuid,
		ServiceName: &serviceName,
	}

	totalPrice, err := client.TotalPrice(ctx, req)
	if err != nil {
		log.Printf("function 'TotalPrice' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, totalPrice)
	log.Printf("function 'TotalPrice' into endpoint service finished successfuly with requestID - '%s'", reqID)
}

func AddSubscription(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'AddSubscription' into endpoint service start with requestID - '%s'", reqID)

	var req SubRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("function 'AddSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
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
		return
	}

	ctx.JSON(http.StatusCreated, resp)
	log.Printf("function 'AddSubscription' into endpoint service finished successfuly with requestID - '%s'", reqID)
}

func UpdatedSubscription(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'UpdatedSubscription' into endpoint service start with requestID - '%s'", reqID)

	service_name := ctx.Param("service_name")
	user_uuid := ctx.Param("user_uuid")
	var req PropsForUpdateREquest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Printf("function 'UpdatedSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	r := &g.UpdatedSubscriptionRequest{
		ReqID:       reqID,
		UserUuid:    user_uuid,
		ServiceName: service_name,
		StartDate:   req.StartDate,
		StopDate:    req.StopDate,
		Price:       req.Price,
	}

	resp, err := client.UpdatedSubscription(ctx, r)
	if err != nil {
		log.Printf("function 'UpdatedSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusAccepted, resp)
	log.Printf("function 'UpdatedSubscription' into endpoint service finished successfuly with requestID - '%s'", reqID)
}

func DeleteSubscription(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	log.Printf("function 'DeleteSubscription' into endpoint service start with requestID - '%s'", reqID)

	userUuid := ctx.Param("user_uuid")
	serviceName := ctx.Param("service_name")

	r := &g.DeleteSubscriptionRequest{
		ReqID:       reqID,
		UserUuid:    userUuid,
		ServiceName: serviceName,
	}

	resp, err := client.DeleteSubscription(ctx, r)
	if err != nil {
		log.Printf("function 'DeleteSubscription' into endpoint service finished with error, requestID - '%s', error - '%v'", reqID, err)
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	ctx.JSON(http.StatusOK, resp)
	log.Printf("function 'DeleteSubscription' into endpoint service finished successfuly with requestID - '%s'", reqID)
}
