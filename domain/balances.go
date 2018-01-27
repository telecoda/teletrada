package domain

import (
	"context"
	"fmt"
	"sync"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
)

type Balance struct {
	sync.RWMutex
	exchanges.ExchangeBalance
	Symbol         Symbol
	Total          float64
	Value          float64
	LatestUSDPrice float64
	LatestUSDValue float64
}

// func (s *server) GetBalances() []*Balance {
// 	return s.livePortfolio.balances
// }

// GetBalances returns current balances
func (s *server) GetBalances(ctx context.Context, in *proto.BalancesRequest) (*proto.BalancesResponse, error) {

	resp := &proto.BalancesResponse{}

	balances := s.livePortfolio.balances

	resp.Balances = make([]*proto.Balance, len(balances))

	for i, balance := range balances {
		resp.Balances[i] = &proto.Balance{
			Symbol:         string(balance.Symbol.GetType()),
			Exchange:       balance.Exchange,
			Free:           float32(balance.Free),
			Locked:         float32(balance.Locked),
			Total:          float32(balance.Total),
			LatestUSDPrice: float32(balance.LatestUSDPrice),
			LatestUSDValue: float32(balance.LatestUSDValue),
		}
	}

	return resp, nil
}

// reprice - updates latestUSDPrice with currency value based on most recent prices
func (b *Balance) reprice() error {
	b.Lock()
	defer b.Unlock()

	price, err := b.Symbol.GetLatestPriceAs(USDT)
	if err != nil {
		return fmt.Errorf("Failed to update value - %s", err)
	}

	b.LatestUSDPrice = price.Price
	b.LatestUSDValue = b.LatestUSDPrice * b.Total

	return nil
}
