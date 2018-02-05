package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
)

type BalanceAs struct {
	sync.RWMutex
	exchanges.Balance
	Total        float64
	At           time.Time
	As           SymbolType
	Price        float64
	Value        float64
	Price24H     float64
	Value24H     float64
	Change24H    float64
	ChangePct24H float64
}

// GetBalances returns current balances
func (s *server) GetBalances(ctx context.Context, req *proto.BalancesRequest) (*proto.BalancesResponse, error) {

	resp := &proto.BalancesResponse{}

	if err := s.updatePortfolios(); err != nil {
		return nil, fmt.Errorf("failed to update portfolios - %s", err)
	}

	balances := s.livePortfolio.balances

	resp.Balances = make([]*proto.Balance, len(balances))

	var err error
	for i, balance := range balances {
		resp.Balances[i], err = balance.toProto()
		if err != nil {
			return nil, err
		}
	}

	return resp, nil
}
