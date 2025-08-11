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
	res := usecase.ListSubscriptions()
	log.Println(res)
	subsSlc := &grpcGen.Subscription{
		ServiceName: res,
	}

	resp := &grpcGen.ListSubscriptionsResponse{
		ReqId:         req.ReqId,
		Subscriptions: []*grpcGen.Subscription{subsSlc},
	}
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
