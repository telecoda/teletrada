package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/telecoda/teletrada/proto"
)

/*

Steps involved in a simulation

There are 2 types of simulation

1 - historic simulation (replayed against old pricing data)
2- live simulation (executed in realtime by scheduler)

Historic simulation process:-

clone live portfolio
apply Strategy to all coin balances

*/

type simulation struct {
	name       string
	*portfolio // original portfolio

	// Simulation specific stuff here

	// historical simulation attributes
	useHistoricalData bool          // do we use historical data
	simFromTime       *time.Time    // when does data start
	simToTime         *time.Time    // when does data end
	dataFrequency     time.Duration // what frequency do we sample the data (normally captured once per minute)

	useRealtimeData bool
}

func (s *server) NewSimulation(simName string, portfolio *portfolio) (*simulation, error) {

	if _, ok := s.simulations[simName]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", simName)
	}

	if portfolio == nil {
		return nil, fmt.Errorf("Portfolio cannot be nil")
	}

	if portfolio.balances == nil || len(portfolio.balances) == 0 {
		return nil, fmt.Errorf("Cannot use a portfolio with no balances for a simulation")
	}

	sim := &simulation{
		name:      simName,
		portfolio: portfolio,
	}

	symbol := SymbolType("ETH")
	as := SymbolType("BTC")
	sellStrat, err := NewPriceAboveStrategy("sell-eth", symbol, as, 0.0545, 100.00)
	if err != nil {
		return nil, err
	}

	buyStrat, err := NewPriceBelowStrategy("buy-eth", symbol, as, 0.0500, 100.00)
	if err != nil {
		return nil, err
	}

	if _, ok := sim.balances[symbol]; !ok {
		return nil, fmt.Errorf("Cannot create simulation as it does not have a balance for %s", symbol)
	}

	sim.balances[symbol].SellStrategy = sellStrat
	sim.balances[symbol].BuyStrategy = buyStrat

	s.simulations[simName] = sim

	// setup simulation parameters

	return sim, nil
}

func (s *simulation) setBuyStrategy(strategy Strategy) error {
	if strategy == nil {
		return fmt.Errorf("Cannot set buy strategy for simulation %q, strategy cannot be nil", s.name)
	}
	// check symbol
	s.Lock()
	defer s.Unlock()
	if symbol, ok := s.balances[strategy.Symbol()]; !ok {
		return fmt.Errorf("Cannot set buy strategy for simulation %q on symbol %q, not in portfolio", s.name, strategy.Symbol())
	} else {
		symbol.Lock()
		symbol.BuyStrategy = strategy
		symbol.Unlock()
	}
	return nil
}

func (s *simulation) setSellStrategy(strategy Strategy) error {
	if strategy == nil {
		return fmt.Errorf("Cannot set sell strategy for simulation %q, strategy cannot be nil", s.name)
	}
	// check symbol
	s.Lock()
	defer s.Unlock()
	if symbol, ok := s.balances[strategy.Symbol()]; !ok {
		return fmt.Errorf("Cannot set sell strategy for simulation %q on symbol %q, not in portfolio", s.name, strategy.Symbol())
	} else {
		symbol.Lock()
		symbol.SellStrategy = strategy
		symbol.Unlock()
	}
	return nil
}

func (s *simulation) runOverHistory(from time.Time, to time.Time, frequency time.Duration) error {

	// validate params
	if from.IsZero() {
		return fmt.Errorf("From time cannot be zero")
	}
	if to.IsZero() {
		return fmt.Errorf("To time cannot be zero")
	}

	if from.After(to) {
		return fmt.Errorf("From time cannot be after to time")
	}

	if frequency.Seconds() != 0 {
		return fmt.Errorf("Frequency cannot be zero")
	}
	s.Lock()
	defer s.Unlock()

	s.simFromTime = &from
	s.simToTime = &to
	s.dataFrequency = frequency
	s.useHistoricalData = true

	return nil

}

func (s *simulation) run() error {
	return nil
}

// GetSimulations returns current simulations
func (s *server) GetSimulations(ctx context.Context, req *proto.GetSimulationsRequest) (*proto.GetSimulationsResponse, error) {

	resp := &proto.GetSimulationsResponse{}

	resp.Simulations = make([]*proto.Simulation, len(s.simulations))

	var err error
	i := 0
	for _, simulation := range s.simulations {
		resp.Simulations[i], err = simulation.toProto()
		if err != nil {
			return nil, err
		}
		i++
	}

	return resp, nil
}
