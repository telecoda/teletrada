package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/api.proto

const (
	port = ":50051"
)

type params struct {
	useMock         bool
	port            int
	loadPricesDir   string
	savePrices      bool
	savePricesDir   string
	priceUpdateFreq time.Duration
	verbose         bool
}

func (p *params) setup() {
	flag.BoolVar(&p.useMock, "usemock", false, "Use mock exchange client")
	flag.BoolVar(&p.verbose, "v", false, "Verbose logging")
	flag.StringVar(&p.loadPricesDir, "loadpricesdir", "priceHistory", "Dir to load historic prices from")
	flag.BoolVar(&p.savePrices, "saveprices", false, "Save new prices as files")
	flag.StringVar(&p.savePricesDir, "savepricesdir", "priceHistory", "Dir to save new prices to")
	flag.DurationVar(&p.priceUpdateFreq, "priceupdatefreq", time.Duration(60*time.Second), "Price update frequency")
	flag.IntVar(&p.port, "port", 13370, "Port for server to listen on")
}

func main() {
	p := &params{}
	p.setup()
	flag.Parse()

	config := domain.Config{
		UseMock:       p.useMock,
		LoadPricesDir: p.loadPricesDir,
		SavePricesDir: p.savePricesDir,
		SavePrices:    p.savePrices,
		UpdateFreq:    p.priceUpdateFreq,
		Verbose:       p.verbose,
	}

	server, err := domain.NewTradaServer(config)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
		return
	}

	if err = server.Init(); err != nil {
		log.Fatalf("Failed to init server - %s", err)
	}

	fmt.Printf("Starting gRPC server on: %s\n", port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	proto.RegisterTeletradaServer(s, server)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
