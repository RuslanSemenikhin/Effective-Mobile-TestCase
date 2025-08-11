package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	grpcGen "github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/rpc/grpc/gen"
	"github.com/RuslanSemenikhin/Effective-Mobile-TestCase.git/usecase"
	"google.golang.org/grpc"
)

type Server struct {
	grpcGen.SubscriptionServiceServer
}

func (s *Server) ListSubscriptions(
	ctx context.Context,
	req *grpcGen.ListSubscriptionsRequest,
) (*grpcGen.ListSubscriptionsResponse, error) {
	log.Printf("start method 'ListSubscriptions' into controller/grpc with requestId - '%s'", req.ReqId)
	res := usecase.ListSubscriptions()
	log.Println(res)
	subsSlc := &grpcGen.Subscription{
		ServiceName: res,
	}

	resp := &grpcGen.ListSubscriptionsResponse{
		ReqId:         req.ReqId,
		Subscriptions: []*grpcGen.Subscription{subsSlc},
	}
	log.Printf("finished successfuly method 'ListSubscriptions' into controller/grpc with requestId - '%s'", req.ReqId)
	return resp, nil
}

func (s *Server) AddSubscription(
	ctx context.Context,
	req *grpcGen.AddSubscriptionRequest,
) (*grpcGen.AddSubscriptionResponse, error) {
	log.Printf("start method 'AddSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	log.Println("===>>>", req.Subscription)
	res := usecase.AddSubscription()
	log.Println(res)
	sub := &grpcGen.Subscription{
		ServiceName: res,
	}

	resp := &grpcGen.AddSubscriptionResponse{
		ReqID:        req.ReqID,
		Subscription: sub,
	}
	log.Printf("finished successfuly method 'AddSubscription' into controller/grpc with requestId - '%s'", req.ReqID)
	return resp, nil
}

func Start(port int) {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	server := grpc.NewServer()
	grpcGen.RegisterSubscriptionServiceServer(server, &Server{})
	log.Printf("controller start on port - '%s'", address)

	if err := server.Serve(listener); err != nil {
		log.Fatalf("Failed to serve controller: %v", err)
	}
}

func Stop() {

}
