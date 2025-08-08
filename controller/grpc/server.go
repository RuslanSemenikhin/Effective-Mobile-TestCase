package grpc

import (
	"fmt"
	"log"
	"net"
	// "google.golang.org/grpc"
)

func Start(port int) {
	address := fmt.Sprintf(":%d", port)
	_, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// server := grpc.NewServer()
}

func Stop() {

}
