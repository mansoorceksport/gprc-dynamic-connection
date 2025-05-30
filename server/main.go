package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	pb "github.com/mansoorceksport/gprc-dynamic-connection/pb"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedEchoServer
}

func (s *server) SayHello(ctx context.Context, req *pb.EchoRequest) (*pb.EchoReply, error) {
	podName := os.Getenv("POD_NAME")
	log.Printf("[B] Received request on pod: %s", podName)
	return &pb.EchoReply{Message: fmt.Sprintf("Hello from %s", podName)}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterEchoServer(s, &server{})
	log.Println("gRPC server listening on :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
