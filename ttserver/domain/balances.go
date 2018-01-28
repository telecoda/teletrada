package domain

import (
	"context"
	"sync"
	"time"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
)

type Balance struct {
	sync.RWMutex
	exchanges.ExchangeBalance
	Total float64
	At    time.Time
}

// GetBalances returns current balances
func (s *server) GetBalances(ctx context.Context, req *proto.BalancesRequest) (*proto.BalancesResponse, error) {

	resp := &proto.BalancesResponse{}

	balances := s.livePortfolio.balances

	resp.Balances = make([]*proto.Balance, len(balances))

	var err error
	for i, balance := range balances {
		resp.Balances[i], err = balance.toProto()
		if err != nil {
			return nil, err
		}

		resp.Balances[i].As = req.As

		// find latest price for trading pair
		price, err := DefaultArchive.GetLatestPriceAs(SymbolType(balance.Symbol), SymbolType(req.As))
		if err != nil {
			return nil, err
		}

		// reprice balance
		resp.Balances[i].AsPrice = float32(price.Price)
		resp.Balances[i].AsValue = float32(price.Price * balance.Total)

	}

	return resp, nil
}
