package domain

import (
	"fmt"
	"log"
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
	livePortfolio *portfolio             // This represents the real live portfolio on the exchange
	simulations   map[string]*simulation // These represent alternate simulated portfolios and their total values
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
		DefaultMetrics, err = newMockMetricsClient("MockMetricsDB")
		if err != nil {
			return nil, err
		}

	} else {
		DefaultClient, err = exchanges.NewBinanceClient()
		if err != nil {
			return nil, err
		}
		DefaultMetrics, err = newMetricsClient(config.InfluxDBName)
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

	s.startTime = time.Now().UTC()

	// scheduler will do a price update immediately
	s.startScheduler()

	if err := s.initPortfolios(); err != nil {
		log.Printf("Failed to initialise portfolio: %s\n", err)
		s.log(fmt.Sprintf("Failed to initialise portfolio: %s", err))
	}

	// TEMP code create simulation
	if _, err := s.NewSimulation("test-sim"); err != nil {
		log.Printf("Failed to create simulation: %s\n", err)
		s.log(fmt.Sprintf("Failed to create simulation: %s", err))
	}

	return nil

}

func (s *server) isVerbose() bool {
	return s.config.Verbose
}
