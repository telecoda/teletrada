package domain

import (
	"fmt"
	"sync"
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

	startScheduler()
	stopScheduler()
	// status logging
	log(msg string)
}

type server struct {
	sync.RWMutex
	livePortfolio *portfolio   // This represents the real live portfolio on the exchange
	simPorts      []*portfolio // These represent alternate simulated portfolios and their total values
	config        Config

	// logging
	statusLog []LogEntry

	// status
	startTime time.Time

	// scheduling
	stopUpdate chan bool
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
		config:     config,
		startTime:  time.Now(),
		stopUpdate: make(chan bool),
	}

	return server, nil
}

func (s *server) Init() error {
	s.Lock()
	defer s.Unlock()

	var err error
	// init DefaultMetrics client
	DefaultMetrics, err = newMetricsClient(s.config.InfluxDBName)
	if err != nil {
		return err
	}

	s.startTime = time.Now().UTC()

	s.startScheduler()

	// update prices immediately
	if err := DefaultArchive.UpdatePrices(); err != nil {
		return fmt.Errorf("Failed to update latest prices: %s", err)
	}

	if err := s.initPortfolios(); err != nil {
		return fmt.Errorf("Failed to initialise portfolio: %s", err)
	}

	return nil

}

func (s *server) isVerbose() bool {
	return s.config.Verbose
}
