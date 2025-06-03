package main

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/mansoorceksport/gprc-dynamic-connection/pb"
	"google.golang.org/grpc"
)

const (
	serviceDNS = "b-service.default.svc.cluster.local"
	port       = 50051
)

var (
	lastHash = ""
	conn     *grpc.ClientConn
	client   pb.EchoClient
)

func hashAddresses(addresses []string) string {
	j, _ := json.Marshal(addresses)
	h := sha256.Sum256(j)
	return fmt.Sprintf("%x", h[:])
}

func resolvePods() ([]string, error) {
	return net.LookupHost(serviceDNS)
}

func connect() {
	address := fmt.Sprintf("dns:///%s:%d", serviceDNS, port)
	var err error
	conn, err = grpc.Dial(
		address,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("Failed to dial gRPC: %v", err)
	}
	client = pb.NewEchoClient(conn)
	log.Println("[INFO] gRPC connected")
}

func disconnect() {
	if conn != nil {
		conn.Close()
		log.Println("[INFO] gRPC connection closed")
	}
}

func main() {
	connect()

	for {
		time.Sleep(2 * time.Second)

		addrs, err := resolvePods()
		if err != nil {
			log.Printf("[ERROR] DNS lookup failed: %v", err)
			continue
		}

		newHash := hashAddresses(addrs)
		if newHash != lastHash {
			log.Printf("[INFO] Detected backend pod change: %v", addrs)
			lastHash = newHash
			disconnect()
			connect()
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		resp, err := client.SayHello(ctx, &pb.EchoRequest{Name: "Mansoor"})
		cancel()

		if err != nil {
			log.Printf("[ERROR] gRPC call failed: %v", err)
			continue
		}
		log.Printf("[INFO] Reply: %s", resp.Message)
	}
}
