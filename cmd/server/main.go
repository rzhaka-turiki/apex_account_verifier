package main

import (
	"log"
	"net"
	"os"

	"github.com/rzhaka-turiki/apex-account-verifier/internal/api"
	grpcserver "github.com/rzhaka-turiki/apex-account-verifier/internal/grpc"
	"github.com/rzhaka-turiki/apex-account-verifier/internal/service"
	apexpb "github.com/rzhaka-turiki/apex-account-verifier/proto/apexpb"

	"google.golang.org/grpc"
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
	verifier := service.NewVerifier(apiClient)
	// handler
	server := grpcserver.NewServer(verifier)
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	apexpb.RegisterApexVerifierServer(grpcServer, server)
	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("grpc server: %v", err)
	}
}
