package grpc

import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "streamHive/streamHive/internal/delivery/grpc/pb"
)

func ServeWithInstance(port string, srv pb.StreamServiceServer) (*grpc.Server, net.Listener) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterStreamServiceServer(grpcServer, srv)

	log.Printf("gRPC server initialized on %s", port)
	return grpcServer, listener
}
