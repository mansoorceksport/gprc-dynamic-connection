package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	pb "github.com/mansoorceksport/gprc-dynamic-connection/pb"
	"google.golang.org/grpc"
)

const (
	serviceDNS = "b-service.default.svc.cluster.local"
	port       = 50051
)

type gRPCManager struct {
	mu     sync.Mutex
	conn   *grpc.ClientConn
	client pb.EchoClient
}

func (g *gRPCManager) connect() {
	address := fmt.Sprintf("dns:///%s:%d", serviceDNS, port)
	conn, err := grpc.Dial(
		address,
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatalf("[gRPCManager] Failed to dial: %v", err)
	}
	g.mu.Lock()
	g.conn = conn
	g.client = pb.NewEchoClient(conn)
	g.mu.Unlock()
	log.Println("[gRPCManager] Connected")
}

func (g *gRPCManager) disconnect() {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.conn != nil {
		_ = g.conn.Close()
		log.Println("[gRPCManager] Connection closed")
	}
}

func (g *gRPCManager) request() {
	g.mu.Lock()
	client := g.client
	g.mu.Unlock()

	if client == nil {
		log.Println("[gRPCManager] Client not initialized")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := client.SayHello(ctx, &pb.EchoRequest{Name: "Mansoor"})
	if err != nil {
		log.Printf("[gRPCManager] Request error: %v", err)
		return
	}
	log.Printf("[gRPCManager] Response: %s", resp.Message)
}

//func hashAddresses(addrs []string) string {
//	j, _ := json.Marshal(addrs)
//	h := sha256.Sum256(j)
//	return fmt.Sprintf("%x", h[:])
//}

func resolverWatcher(interval time.Duration, changed chan<- struct{}) {
	var lastHash string

	for {
		time.Sleep(interval)
		addrs, err := net.LookupHost(serviceDNS)
		if err != nil {
			log.Printf("[resolverWatcher] DNS error: %v", err)
			continue
		}

		currentHash := hashAddresses(addrs)
		if currentHash != lastHash {
			log.Printf("[resolverWatcher] Pod list changed: %v", addrs)
			lastHash = currentHash
			changed <- struct{}{}
		}
	}
}

func main() {
	log.Println("[Main] Starting gRPC client with pod change detection...")

	changeSignal := make(chan struct{})
	manager := &gRPCManager{}
	manager.connect()

	go resolverWatcher(10*time.Second, changeSignal)

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			manager.request()

		case <-changeSignal:
			log.Println("[Main] Triggered re-dial due to backend pod change")
			manager.disconnect()
			manager.connect()
		}
	}
}
