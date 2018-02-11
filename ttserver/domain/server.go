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
	log(msg string)
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
	InfluxDBName   string
	InfluxUsername string
	InfluxPassword string
	UpdateFreq     time.Duration
	Verbose        bool
	Port           int
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

	return server, nil
}

func (s *server) Init() error {

	s.startTime = time.Now().UTC()

	// if s.config.LoadPricesDir != "" {
	// 	go func() {
	// 		s.log("Loading historic prices from filesystem")
	// 		if err := DefaultArchive.LoadPrices(s.config.LoadPricesDir); err != nil {
	// 			s.log(fmt.Sprintf("ERROR:Failed to load historic prices: %s", err))
	// 		}
	// 	}()
	// }

	DefaultArchive.StartUpdater(s.config.UpdateFreq)

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
