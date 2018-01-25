package domain

import (
	"fmt"
	"sync"

	"github.com/telecoda/teletrada/exchanges"
)

type balance struct {
	sync.RWMutex
	exchanges.ExchangeBalance
	symbol         Symbol
	Total          float64
	Value          float64
	LatestUSDPrice float64
	LatestUSDValue float64
}

// reprice - updates latestUSDPrice with currency value based on most recent prices
func (b *balance) reprice() error {
	b.Lock()
	defer b.Unlock()

	price, err := b.symbol.GetLatestPriceAs(USDT)
	if err != nil {
		return fmt.Errorf("Failed to update value - %s", err)
	}

	b.LatestUSDPrice = price.Price
	b.LatestUSDPrice = b.LatestUSDPrice * b.Total
	return nil
}
