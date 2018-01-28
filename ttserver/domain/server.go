package domain

import (
	"fmt"
	"time"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
)

const (
	DATE_FORMAT = "2006-01-02 03:04:05"
)

var DefaultClient exchanges.ExchangeClient

type Server interface {
	proto.TeletradaServer
	Init() error

	// status logging
	//GetLog() []Log
	log(msg string)

	// balances
	ListBalances()
	//GetBalances() []*Balance
}

type server struct {
	livePortfolio *portfolio   // This represents the real live portfolio on the exchange
	simPorts      []*portfolio // These represent alternate simulated portfolios and their total values
	config        Config

	// logging
	statusLog []LogEntry
	startTime time.Time
}

type Config struct {
	UseMock        bool
	LoadPricesDir  string
	SavePricesDir  string
	SavePrices     bool
	UpdateDuration time.Duration
	Verbose        bool
}

func NewTradaServer(config Config) (Server, error) {

	var err error
	if config.UseMock {
		DefaultClient, err = exchanges.NewMockClient()
		if err != nil {
			return nil, err
		}
	} else {
		DefaultClient, err = exchanges.NewBinanceClient()
		if err != nil {
			return nil, err
		}
	}

	server := &server{
		config: config,
	}

	if err := server.Init(); err != nil {
		return nil, err
	}

	server.ListBalances()
	return server, nil
}

func (s *server) Init() error {

	s.startTime = time.Now().UTC()

	if s.config.LoadPricesDir != "" {
		s.log("Loading historic prices from filesystem")
		if err := DefaultArchive.LoadPrices(s.config.LoadPricesDir); err != nil {
			return fmt.Errorf("Failed to load historic prices: %s", err)
		}
	}

	if s.config.SavePrices {
		s.log("Starting price persistence")
		if err := DefaultArchive.StartPersistence(s.config.SavePricesDir); err != nil {
			return err
		}
	}

	if err := DefaultArchive.UpdatePrices(); err != nil {
		return fmt.Errorf("Failed to update latest prices: %s", err)
	}

	if err := s.initPortfolios(); err != nil {
		return fmt.Errorf("Failed to initialise portfolio: %s", err)
	}

	if err := s.updatePortfolios(); err != nil {
		return fmt.Errorf("Failed to update portfolio: %s", err)
	}

	return nil

}

func (s *server) isVerbose() bool {
	return s.config.Verbose
}

func (s *server) ListBalances() {
	fmt.Printf("Live Portfolio\n")
	fmt.Printf("==============\n")
	s.livePortfolio.ListBalances()

	fmt.Printf("Simulated Portfolios\n")
	fmt.Printf("====================\n")
	for _, p := range s.simPorts {
		p.ListBalances()
	}
}
