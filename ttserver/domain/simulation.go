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
	id          string
	name        string
	isRunning   bool
	startedTime *time.Time // time simulation started
	stoppedTime *time.Time // time simulation stopped

	*portfolio // original portfolio

	// Simulation specific stuff here

	// historical simulation attributes
	useHistoricalData bool          // do we use historical data
	simFromTime       *time.Time    // when does data start
	simToTime         *time.Time    // when does data end
	dataFrequency     time.Duration // what frequency do we sample the data (normally captured once per minute)

	useRealtimeData bool
}

func (s *server) getSimulation(id string) (*simulation, error) {
	// lookup simulation
	s.RLock()
	defer s.RUnlock()

	if sim, ok := s.simulations[id]; !ok {
		return nil, fmt.Errorf("Simulation Id: %s not found", id)
	} else {
		return sim, nil
	}
}

func (s *server) setSimulation(sim *simulation) {
	// lookup simulation
	s.Lock()
	defer s.Unlock()

	s.simulations[sim.id] = sim
}

// CreateSimulation creates a new simulation
func (s *server) CreateSimulation(ctx context.Context, in *proto.CreateSimulationRequest) (*proto.CreateSimulationResponse, error) {
	resp := &proto.CreateSimulationResponse{}

	return resp, nil
}

// StartSimulation starts a simulation running
func (s *server) StartSimulation(ctx context.Context, req *proto.StartSimulationRequest) (*proto.StartSimulationResponse, error) {

	if req.Id == "" {
		return nil, fmt.Errorf("You must provide a simulation Id")
	}

	sim, err := s.getSimulation(req.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get simulation - %s", err)
	}

	sim.Lock()
	defer sim.Unlock()

	if sim.isRunning {
		return nil, fmt.Errorf("Simulation Id: %s is already started", req.Id)
	}

	sim.isRunning = true
	now := time.Now().UTC()
	sim.startedTime = &now

	s.setSimulation(sim)

	s.log(fmt.Sprintf("Simulation: %s started running", sim.id))
	go sim.run()

	resp := &proto.StartSimulationResponse{}

	return resp, nil
}

// StopSimulation stops a simulation running
func (s *server) StopSimulation(ctx context.Context, req *proto.StopSimulationRequest) (*proto.StopSimulationResponse, error) {
	resp := &proto.StopSimulationResponse{}

	if req.Id == "" {
		return nil, fmt.Errorf("You must provide a simulation Id")
	}

	sim, err := s.getSimulation(req.Id)
	if err != nil {
		return nil, fmt.Errorf("Failed to get simulation - %s", err)
	}

	sim.Lock()
	defer sim.Unlock()

	if !sim.isRunning {
		return nil, fmt.Errorf("Simulation Id: %s is not running", req.Id)
	}

	sim.isRunning = false
	now := time.Now().UTC()
	sim.stoppedTime = &now

	s.setSimulation(sim)

	s.log(fmt.Sprintf("Simulation: %s stop requested", sim.id))

	return resp, nil
}

func (s *server) NewSimulation(id, simName string, portfolio *portfolio) (*simulation, error) {

	if _, ok := s.simulations[id]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", id)
	}

	if portfolio == nil {
		return nil, fmt.Errorf("Portfolio cannot be nil")
	}

	if portfolio.balances == nil || len(portfolio.balances) == 0 {
		return nil, fmt.Errorf("Cannot use a portfolio with no balances for a simulation")
	}

	sim := &simulation{
		id:        id,
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

	s.simulations[id] = sim

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

func (s *simulation) run() {
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
