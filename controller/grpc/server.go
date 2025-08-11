package grpc

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/internal/models"
	grpcGen "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/usecase"
	"google.golang.org/grpc"
)

type Server struct {
	grpcGen.SubscriptionServiceServer
	db *sql.DB
}

func (s *Server) ListSubscriptions(
	ctx context.Context,
	req *grpcGen.ListSubscriptionsRequest,
) (*grpcGen.ListSubscriptionsResponse, error) {
	log.Printf("start method 'ListSubscriptions' into controller/grpc with requestId - '%s'", req.ReqId)
	res, err := usecase.ListSubscriptions(
		ctx,
		s.db,
		req.ServiceName,
		req.UserUuid,
		req.StartDate,
		req.StopDate,
	)
	if err != nil {
		return nil, err
	}

	subsInfos := []*grpcGen.ListSubscriptionsInfo{}

	for _, r := range res {
		subInfo := &grpcGen.ListSubscriptionsInfo{
			UserUuid:    r.UserUuid,
			ServiceName: r.ServiceName.String,
			Price:       r.ServicePrice,
		}
		subsInfos = append(subsInfos, subInfo)
	}

	subsSlc := &grpcGen.ListSubscriptionsResponse{
		ReqId:         req.ReqId,
		Subscriptions: subsInfos,
	}
	log.Printf("finished successfuly method 'ListSubscriptions' into controller/grpc with requestId - '%s'", req.ReqId)
	return subsSlc, nil
}

func (s *Server) AddSubscription(
	ctx context.Context,
	req *grpcGen.AddSubscriptionRequest,
) (*grpcGen.AddSubscriptionResponse, error) {
	log.Printf("start method 'AddSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	res, err := usecase.AddSubscription(
		ctx,
		s.db,
		req.Subscription.ServiceName,
		req.Subscription.UserUuid,
		req.Subscription.Price,
		req.Subscription.StartDate,
		req.Subscription.StopDate,
	)
	if err != nil {
		return nil, err
	}

	sub := &grpcGen.Subscription{
		ServiceName: req.Subscription.ServiceName,
		Price:       req.Subscription.Price,
		UserUuid:    res.UserUuid,
		StartDate:   res.StartDate.Format("01.2006"),
		StopDate:    res.StopDate.Format("01.2006"),
	}

	resp := &grpcGen.AddSubscriptionResponse{
		ReqID:        req.ReqID,
		Subscription: sub,
	}
	log.Printf("finished successfuly method 'AddSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	return resp, nil
}

func (s *Server) UpdatedSubscription(
	ctx context.Context,
	req *grpcGen.UpdatedSubscriptionRequest,
) (*grpcGen.UpdatedSubscriptionResponse, error) {
	log.Printf("start method 'UpdatedSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	res := usecase.UpdatedSubscription(
		ctx,
		s.db,
		req.ServiceName,
		req.UserUuid,
		req.Price,
		req.StartDate,
		req.StopDate,
	)
	log.Println(res)
	sub := &grpcGen.Subscription{
		ServiceName: res,
	}

	resp := &grpcGen.UpdatedSubscriptionResponse{
		ReqID:        req.ReqID,
		Subscription: sub,
	}
	log.Printf("finished successfuly method 'AddSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	return resp, nil
}

func (s *Server) DeleteSubscription(
	ctx context.Context,
	req *grpcGen.DeleteSubscriptionRequest,
) (*grpcGen.DeleteSubscriptionResponse, error) {
	log.Printf("start method 'DeleteSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	res := usecase.DeleteSubscription(
		ctx,
		s.db,
		req.ServiceName,
		req.UserUuid,
	)
	log.Println(res)
	sub := &grpcGen.Subscription{
		ServiceName: res,
	}

	resp := &grpcGen.DeleteSubscriptionResponse{
		ReqID:        req.ReqID,
		Subscription: sub,
	}
	log.Printf("finished successfuly method 'DeleteSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	return resp, nil
}

func Start(config *models.Config, dbCon *sql.DB) {
	address := fmt.Sprintf(":%d", config.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	grpcGen.RegisterSubscriptionServiceServer(server, &Server{
		db: dbCon,
	})
	log.Printf("controller start on port - '%s'", address)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve controller: %v", err)
	}
}

func Stop() {

}
