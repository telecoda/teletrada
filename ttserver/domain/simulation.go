package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/servertime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	// simulation acts and behaves the same as a "normal" portfolio
	// this is held in the *portfolio variable
	*portfolio // current simulated

	// realNow is a reference to the real portfolio now
	// so any difference between the simulated value and real value can be compared
	realNow *portfolio

	// realStart is a cloned copy of the real portfolio
	// at the beginning of the simulation
	realAtStart *portfolio

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
func (s *server) CreateSimulation(ctx context.Context, req *proto.CreateSimulationRequest) (*proto.CreateSimulationResponse, error) {
	// validate request
	if req.Id == "" {
		// if no id is provided generate one
		req.Id = randSeq(10)
	}

	if req.Name == "" {
		// if no name is provided generate one
		req.Name = "Created simulation"
	}

	sim, err := s.newSimulation(req.Id, req.Name)
	if err != nil {
		return nil, err
	}

	pSim, err := sim.toProto()
	if err != nil {
		return nil, err
	}

	resp := &proto.CreateSimulationResponse{
		Simulation: pSim,
	}

	return resp, nil
}

// StartSimulation starts a simulation running
func (s *server) StartSimulation(ctx context.Context, req *proto.StartSimulationRequest) (*proto.StartSimulationResponse, error) {

	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "You must provide a simulation Id")
	}

	sim, err := s.getSimulation(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to get simulation - %s", err)
	}

	sim.Lock()
	defer sim.Unlock()

	if sim.isRunning {
		return nil, status.Errorf(codes.Unavailable, "Simulation Id: %s is already started", req.Id)
	}

	// validate when
	if _, ok := proto.StartSimulationRequestWhenOptions_name[int32(req.When)]; !ok {
		return nil, status.Errorf(codes.InvalidArgument, "When value %d is not valid", req.When)
	}

	now := servertime.Now()

	switch req.When {
	case proto.StartSimulationRequest_LAST_DAY:
		sim.simToTime = &now
		from := now.AddDate(0, 0, -1)
		sim.simFromTime = &from
		sim.useHistoricalData = true
	case proto.StartSimulationRequest_LAST_WEEK:
		sim.simToTime = &now
		from := now.AddDate(0, 0, -7)
		sim.simFromTime = &from
		sim.useHistoricalData = true
	case proto.StartSimulationRequest_LAST_MONTH:
		sim.simToTime = &now
		from := now.AddDate(0, 0, -30)
		sim.simFromTime = &from
		sim.useHistoricalData = true
	case proto.StartSimulationRequest_THE_LOT:
		sim.simToTime = &now
		from := now.AddDate(-10, 0, 0) // 10 year should be long enough..
		sim.simFromTime = &from
		sim.useHistoricalData = true
	case proto.StartSimulationRequest_NOW_REALTIME:
		sim.useRealtimeData = true
	default:
		return nil, status.Errorf(codes.InvalidArgument, "When value %d is not valid", req.When)
	}

	if req.When == proto.StartSimulationRequest_NOW_REALTIME {
		sim.useRealtimeData = true
		sim.useHistoricalData = false
	}

	if !sim.useHistoricalData && !sim.useRealtimeData {
		return nil, status.Errorf(codes.Unavailable, "Must enable either historical or realtime data")
	}

	if sim.useHistoricalData && sim.useRealtimeData {
		return nil, status.Errorf(codes.Unavailable, "Simulation cannot be run in historical and realtime mode simultaneously")
	}

	// make a copy of the real portfolio before starting
	// so we can use it to compare results against

	realAtStart, err := sim.realNow.clone()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to clone portfolio - %s", err)
	}

	sim.realAtStart = realAtStart

	go sim.run()

	resp := &proto.StartSimulationResponse{}

	return resp, nil
}

// StopSimulation stops a simulation running
func (s *server) StopSimulation(ctx context.Context, req *proto.StopSimulationRequest) (*proto.StopSimulationResponse, error) {
	resp := &proto.StopSimulationResponse{}

	if req.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "You must provide a simulation Id")
	}

	sim, err := s.getSimulation(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to get simulation - %s", err)
	}

	sim.Lock()
	defer sim.Unlock()

	if !sim.isRunning {
		return nil, status.Errorf(codes.Unavailable, "Simulation Id: %s is not running", req.Id)
	}

	sim.isRunning = false
	now := servertime.Now()
	sim.stoppedTime = &now

	s.setSimulation(sim)

	DefaultLogger.log(fmt.Sprintf("Simulation: %s stop requested", sim.id))

	return resp, nil
}

func (s *server) newSimulation(id, simName string) (*simulation, error) {

	if _, ok := s.simulations[id]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", id)
	}

	real := s.livePortfolio

	if real == nil {
		return nil, fmt.Errorf("Portfolio cannot be nil")
	}

	if real.balances == nil || len(real.balances) == 0 {
		return nil, fmt.Errorf("Cannot use a portfolio with no balances for a simulation")
	}

	clonedPort, err := real.clone()
	if err != nil {
		return nil, fmt.Errorf("Failed to clone real portfolio: %s", err)
	}

	sim := &simulation{
		id:        id,
		name:      simName,
		portfolio: clonedPort,
		realNow:   real,
	}

	s.simulations[id] = sim

	return sim, nil
}

func (s *simulation) SetBuyStrategy(strategy Strategy) error {
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

func (s *simulation) SetSellStrategy(strategy Strategy) error {
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

func (s *simulation) run() {

	s.Lock()
	s.isRunning = true
	now := servertime.Now()
	s.startedTime = &now
	// take a copy of portfolio at start
	s.Unlock()

	// sleep a little at the start
	// just to help the tests do little check...
	time.Sleep(500 * time.Millisecond)

	defer func() {
		s.Lock()
		s.isRunning = false
		now := servertime.Now()
		s.stoppedTime = &now
		s.Unlock()
	}()

	if s.useHistoricalData {

		frequency := time.Duration(5 * time.Minute)

		err := s.runOverHistory(frequency)
		if err != nil {
			DefaultLogger.log(fmt.Sprintf("Error running historic simulation: %s - %s", s.id, err))
			return
		}

	}

	if s.useRealtimeData {
		err := s.runRealtime()
		if err != nil {
			DefaultLogger.log(fmt.Sprintf("Error running realtime simulation: %s - %s", s.id, err))
			return
		}

	}
}

func (s *simulation) runOverHistory(frequency time.Duration) error {

	// validate params
	if s.simFromTime == nil {
		return fmt.Errorf("From time cannot be nil")
	}
	if s.simToTime == nil {
		return fmt.Errorf("To time cannot be nil")
	}
	if s.simFromTime.IsZero() {
		return fmt.Errorf("From time cannot be zero")
	}
	if s.simToTime.IsZero() {
		return fmt.Errorf("To time cannot be zero")
	}

	if s.simFromTime.After(*s.simToTime) {
		return fmt.Errorf("From time cannot be after to time")
	}

	if frequency.Seconds() == 0 {
		return fmt.Errorf("Frequency cannot be zero")
	}
	s.Lock()
	defer s.Unlock()

	s.dataFrequency = frequency
	s.useHistoricalData = true

	DefaultLogger.log(fmt.Sprintf("Historical simulation: %s started", s.id))

	// Save a before version of the portfolio
	before, err := s.portfolio.clone()
	if err != nil {
		return fmt.Errorf("Error cloning portfolio during historical simulation: %s - %s", s.id, err)

	}
	// Replay all prices between dates
	toTime := *s.simToTime

	for priceTime := *s.simFromTime; priceTime.Before(toTime) || priceTime.Equal(toTime); priceTime = priceTime.Add(s.dataFrequency) {
		// reprice current portfolio at this time
		if err := s.portfolio.repriceAt(priceTime); err != nil {
			return fmt.Errorf("Error repriced simulated portfolio at: %s - %s", priceTime.String(), err)
		}

		// now coins have correct price for time
		// execute strategies

		for symbol, balance := range s.portfolio.balances {
			if balance.SellStrategy != nil {
				// exec Sell strat
				sell, err := balance.SellStrategy.ConditionMet(priceTime)
				if err != nil {
					return fmt.Errorf("Error executing sell strategy for symbol: %s - %s", symbol, err)
				}
				if sell {
					// sell, Sell, SELL!
				}
			}
			if balance.BuyStrategy != nil {
				// exec Buy strat
			}
		}

		// DefaultLogger.log(fmt.Sprintf("Reading prices for %s", priceTime))
		// process all symbols in portfolio
		//for s.portfolio

	}

	// Compare portfolio afterwards
	diff, err := s.portfolio.diff(before)
	if err != nil {
		return fmt.Errorf("Error comparing portfolio differences - %s", err)
	}

	// print portfolio diffs
	if true == false {
		// TEMP: skip for now
		diff.print()
	}

	DefaultLogger.log(fmt.Sprintf("Historical simulation: %s ended", s.id))
	return nil
}

func (s *simulation) runRealtime() error {
	DefaultLogger.log(fmt.Sprintf("Realtime simulation: %s started", s.id))

	DefaultLogger.log(fmt.Sprintf("Realtime simulation: %s ended", s.id))
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
