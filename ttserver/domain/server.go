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
}

type server struct {
	sync.RWMutex
	livePortfolio *portfolio             // This represents the real live portfolio on the exchange
	simulations   map[string]*simulation // These represent alternate simulated portfolios and their total values
	config        Config

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

	DefaultLogger = NewLogger(config.Verbose)

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
		startTime:  ServerTime(),
		stopUpdate: make(chan bool),
	}

	return server, nil
}

func (s *server) Init() error {
	s.Lock()
	defer s.Unlock()

	s.startTime = ServerTime()

	// scheduler will do a price update immediately
	s.startScheduler()

	if err := s.initPortfolios(); err != nil {
		DefaultLogger.log(fmt.Sprintf("Failed to initialise portfolio: %s", err))
	}

	// // TEMP code create simulation

	// testSim, err := s.NewSimulation("dummy-init-sim-id", "dummy-init-sim", s.livePortfolio)
	// if err != nil {
	// 	DefaultLogger.log(fmt.Sprintf("Failed to create simulation: %s", err))
	// 	return nil
	// }

	// // Set Buy/Sell strategies on some symbols
	// ethBuy, err := NewPriceAboveStrategy("buy-eth", SymbolType("ETH"), SymbolType("USDT"), 20.00, 100.0)
	// if err != nil {
	// 	DefaultLogger.log(fmt.Sprintf("Failed to create buy strategy: %s", err))
	// 	return nil
	// }
	// ethSell, err := NewPriceBelowStrategy("sell-eth", SymbolType("ETH"), SymbolType("USDT"), 10.00, 100.0)
	// if err != nil {
	// 	DefaultLogger.log(fmt.Sprintf("Failed to create sell strategy: %s", err))
	// 	return nil
	// }

	// if err := testSim.SetBuyStrategy(ethBuy); err != nil {
	// 	DefaultLogger.log(fmt.Sprintf("Failed to set buy strategy: %s", err))
	// 	return nil
	// }

	// if err := testSim.SetSellStrategy(ethSell); err != nil {
	// 	DefaultLogger.log(fmt.Sprintf("Failed to set buy strategy: %s", err))
	// 	return nil
	// }

	return nil

}
