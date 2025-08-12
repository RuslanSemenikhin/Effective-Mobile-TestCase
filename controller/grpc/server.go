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
	)
	if err != nil {
		return nil, err
	}

	resp := []*grpcGen.Subscription{}
	for _, r := range res {
		resp = append(resp, &grpcGen.Subscription{
			ServiceName: r.ServiceName.String,
			Price:       int64(r.ServicePrice.Int32),
			UserUuid:    r.UserUuid,
			StartDate:   r.StartDate.Format("01.2006"),
			StopDate:    r.StopDate.Format("01.2006"),
		})
	}

	r := &grpcGen.ListSubscriptionsResponse{
		ReqId:         req.ReqId,
		Subscriptions: resp,
	}

	log.Printf("finished successfuly method 'ListSubscriptions' into controller/grpc with requestId - '%s'", req.ReqId)
	return r, nil
}

func (s *Server) TotalPrice(
	ctx context.Context,
	req *grpcGen.TotalPriceRequest,
) (*grpcGen.TotalPriceResponse, error) {
	log.Printf("start method 'TotalPrice' into controller/grpc with requestId - '%s'", req.ReqID)
	totalPrice, err := usecase.TotalPrice(
		ctx,
		s.db,
		req.StartDate,
		req.StopDate,
		req.UserUuid,
		req.ServiceName,
	)
	if err != nil {
		return nil, err
	}

	resp := &grpcGen.TotalPriceResponse{
		ReqID:      req.ReqID,
		TotalPrice: totalPrice,
	}
	log.Printf("finished successfuly method 'TotalPrice' into controller/grpc with requestId - '%s'", req.ReqID)
	return resp, nil

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
	err := usecase.DeleteSubscription(
		ctx,
		s.db,
		req.ServiceName,
		req.UserUuid,
	)
	if err != nil {
		return nil, err
	}

	resp := &grpcGen.DeleteSubscriptionResponse{
		ReqID:       req.ReqID,
		UserUuid:    req.UserUuid,
		ServiceName: req.ServiceName,
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
