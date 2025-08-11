package control

import (
	"fmt"
	"net/http"

	g "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListSubscriptions(
	ctx *gin.Context,
	client g.SubscriptionServiceClient,
) {
	reqID := uuid.New().String()
	fmt.Printf("request ID - '%s'", reqID)
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
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, subscriptions)
}
