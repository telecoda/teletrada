package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/telecoda/teletrada/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

type server struct{}

func run() {
	fmt.Printf("Starting gRPC server on: %s\n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterTraderServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// GetBalances returns current balances
func (s *server) GetBalances(ctx context.Context, in *proto.BalancesRequest) (*proto.BalancesResponse, error) {
	return &proto.BalancesResponse{Message: "Balances here"}, nil
}
