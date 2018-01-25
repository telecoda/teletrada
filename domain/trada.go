package domain

import (
	"fmt"
	"time"

	"github.com/telecoda/teletrada/exchanges"
)

const (
	DATE_FORMAT = "2006-01-02 03:04:05"
)

var DefaultClient exchanges.ExchangeClient

type Trada interface {
	Init() error
	ListBalances()
}

type trada struct {
	livePortfolio *portfolio   // This represents the real live portfolio on the exchange
	simPorts      []*portfolio // These represent alternate simulated portfolios and their total values
	config        Config
}

type Config struct {
	UseMock        bool
	LoadPricesDir  string
	SavePricesDir  string
	SavePrices     bool
	UpdateDuration time.Duration
}

func NewTrada(config Config) (Trada, error) {

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

	trada := &trada{
		config: config,
	}

	if err := trada.Init(); err != nil {
		return nil, err
	}

	trada.ListBalances()
	return trada, nil
}

func (t *trada) Init() error {

	if t.config.LoadPricesDir != "" {
		fmt.Printf("Loading historic prices from filesystem\n")
		if err := DefaultArchive.LoadPrices(t.config.LoadPricesDir); err != nil {
			return fmt.Errorf("Failed to load historic prices: %s", err)
		}
	}

	if t.config.SavePrices {
		fmt.Printf("Starting price persistence\n")
		if err := DefaultArchive.StartPersistence(t.config.SavePricesDir); err != nil {
			return err
		}
	}

	if err := DefaultArchive.UpdatePrices(); err != nil {
		return fmt.Errorf("Failed to update latest prices: %s", err)
	}

	if err := t.initPortfolios(); err != nil {
		return fmt.Errorf("Failed to initialise portfolio: %s", err)
	}

	if err := t.updatePortfolios(); err != nil {
		return fmt.Errorf("Failed to update portfolio: %s", err)
	}

	return nil

}

func (t *trada) ListBalances() {
	fmt.Printf("Live Portfolio\n")
	fmt.Printf("==============\n")
	t.livePortfolio.ListBalances()

	fmt.Printf("Simulated Portfolios\n")
	fmt.Printf("====================\n")
	for _, p := range t.simPorts {
		p.ListBalances()
	}
}
