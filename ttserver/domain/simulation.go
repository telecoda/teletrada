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
	*portfolio // original portfolio

	// Simulation specific stuff here

	// historic simulation attributes
	useHistoricData bool          // do we use historic data
	simFromTime     *time.Time    // when does data start
	simToTime       *time.Time    // when does data end
	dataFrequency   time.Duration // what frequency do we sample the data (normally captured once per minute)

	useRealtimeData bool
}

func (s *server) NewSimulation(simName string) (*simulation, error) {

	if _, ok := s.simulations[simName]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", simName)
	}

	clonedPort, err := s.livePortfolio.clone(simName)
	if err != nil {
		return nil, err
	}

	sim := &simulation{
		portfolio: clonedPort,
	}

	symbol := SymbolType("TRX")
	as := SymbolType("USDT")
	sellStrat, err := NewPriceAboveStrategy("sell-trx", symbol, as, 0.0545, 100.00)
	if err != nil {
		return nil, err
	}

	buyStrat, err := NewPriceBelowStrategy("buy-trx", symbol, as, 0.0500, 100.00)
	if err != nil {
		return nil, err
	}

	sim.balances[symbol].SellStrategy = sellStrat
	sim.balances[symbol].SellStrategy = buyStrat

	s.simulations[simName] = sim

	// setup simulation parameters

	return sim, nil
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
