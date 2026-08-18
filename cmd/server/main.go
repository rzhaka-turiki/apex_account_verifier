package main

import (
	"log"
	"net"
	"os"

	"github.com/rzhaka-turiki/apex_account_verifier/internal/api"
	grpcserver "github.com/rzhaka-turiki/apex_account_verifier/internal/grpc"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/queue"
	"github.com/rzhaka-turiki/apex_account_verifier/internal/service"
	apexpb "github.com/rzhaka-turiki/apex_account_verifier/proto/apexpb/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	baseURL := os.Getenv("APEX_API_BASE_URL")
	authToken := os.Getenv("APEX_API_AUTH_TOKEN")
	if baseURL == "" {
		log.Fatal("APEX_API_BASE_URL is not set")
	}
	if authToken == "" {
		log.Fatal("APEX_API_AUTH_TOKEN is not set")
	}
	apiClient := api.NewClient(baseURL, authToken)
	reqQueue := queue.NewQueue(10000, 5, 4)
	verifier := service.NewVerifier(apiClient, reqQueue)
	// handler
	server := grpcserver.NewServer(verifier)
	listener, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	apexpb.RegisterApexVerifierServer(grpcServer, server)
	reflection.Register(grpcServer)
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}
