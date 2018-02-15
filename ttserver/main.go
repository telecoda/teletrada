package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/api.proto

const (
	// env var names
	INFLUX_DB_NAME  = "INFLUX_DB_NAME"
	INFLUX_USERNAME = "INFLUX_USERNAME"
	INFLUX_PASSWORD = "INFLUX_PASSWORD"

	INFLUX_DATABASE      = "teletrada"
	TEST_INFLUX_DATABASE = "testteletrada"
)

type params struct {
	useMock       bool
	port          int
	loadPricesDir string
	updateFreq    time.Duration
	verbose       bool
}

func (p *params) setup() {
	flag.BoolVar(&p.useMock, "usemock", false, "Use mock exchange client")
	flag.BoolVar(&p.verbose, "v", false, "Verbose logging")
	flag.StringVar(&p.loadPricesDir, "loadpricesdir", "priceHistory", "Dir to load historic prices from")
	flag.DurationVar(&p.updateFreq, "updatefreq", time.Duration(60*time.Second), "Update frequency")
	flag.IntVar(&p.port, "port", 13370, "Port for server to listen on")
}

func main() {
	p := &params{}
	p.setup()
	flag.Parse()

	config := domain.Config{
		UseMock:        p.useMock,
		LoadPricesDir:  p.loadPricesDir,
		InfluxDBName:   os.Getenv(INFLUX_DB_NAME),
		InfluxUsername: os.Getenv(INFLUX_USERNAME),
		InfluxPassword: os.Getenv(INFLUX_PASSWORD),
		UpdateFreq:     p.updateFreq,
		Verbose:        p.verbose,
		Port:           p.port,
	}

	// if no env vars, use defaults
	if config.InfluxDBName == "" {
		config.InfluxDBName = INFLUX_DATABASE
	}

	server, err := domain.NewTradaServer(config)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
		return
	}

	if err = server.Init(); err != nil {
		log.Fatalf("Failed to init server - %s", err)
	}

	fmt.Printf("Starting gRPC server on: :%d\n", config.Port)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
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
