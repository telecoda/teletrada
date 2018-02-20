package domain

import "fmt"

/*

Steps involved in a simulation

There are 2 types of simulation

1 - historic simulation (replayed against old pricing data)
2- live simulation (executed in realtime by scheduler)

Historic simulation process:-

clone live portfolio
apply Strategy to all coin balances

*/

func (s *server) NewSimulation(simName string) (*portfolio, error) {

	if _, ok := s.simPorts[simName]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", simName)
	}

	sim, err := s.livePortfolio.clone(simName)
	if err != nil {
		return nil, err
	}

	s.simPorts[simName] = sim

	symbol := SymbolType("TRX")
	as := SymbolType("USDT")
	sellStrat, err := NewPriceAboveStrategy("sim-001", symbol, as, 0.0545, 100.00)
	if err != nil {
		return nil, err
	}

	sim.balances[symbol].SellStrategy = sellStrat

	return sim, nil
}
