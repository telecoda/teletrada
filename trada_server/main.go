package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/telecoda/teletrada/domain"
)

//go:generate protoc -I ../proto --go_out=plugins=grpc:../proto ../proto/api.proto

type params struct {
	useMock         bool
	port            int
	loadPricesDir   string
	savePrices      bool
	savePricesDir   string
	priceUpdateFreq time.Duration
}

func (p *params) setup() {
	flag.BoolVar(&p.useMock, "usemock", false, "Use mock exchange client")
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
	}

	trada, err := domain.NewTrada(config)
	if err != nil {
		fmt.Println(err)
		return
	}

	if err = trada.Init(); err != nil {
		log.Fatalf("Failed to init trada - %s", err)
	}

	// start grpc server
	run()

}
