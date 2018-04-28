package domain

import (
	"context"
	"fmt"
	"time"

	"github.com/telecoda/teletrada/proto"
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

type Simulation interface {
	GetID() string
	GetName() string
	SetBuyStrategy(strategy Strategy) error
	SetSellStrategy(strategy Strategy) error
}

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

	// validate when
	if _, ok := proto.StartSimulationRequestWhenOptions_name[int32(req.When)]; !ok {
		return nil, fmt.Errorf("When value %d is not valid", req.When)
	}

	now := ServerTime()

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
		return nil, fmt.Errorf("When value %d is not valid", req.When)
	}

	if req.When == proto.StartSimulationRequest_NOW_REALTIME {
		sim.useRealtimeData = true
		sim.useHistoricalData = false
	}

	if !sim.useHistoricalData && !sim.useRealtimeData {
		return nil, status.Newf(codes.Unavailable, "Must enable either historical or realtime data").Err()
	}

	if sim.useHistoricalData && sim.useRealtimeData {
		return nil, status.Newf(codes.Unavailable, "Simulation cannot be run in historical and realtime mode simultaneously").Err()
	}

	go sim.run()

	resp := &proto.StartSimulationResponse{}

	return resp, nil
}

// StopSimulation stops a simulation running
func (s *server) StopSimulation(ctx context.Context, req *proto.StopSimulationRequest) (*proto.StopSimulationResponse, error) {
	resp := &proto.StopSimulationResponse{}

	if req.Id == "" {
		return nil, status.Newf(codes.InvalidArgument, "You must provide a simulation Id").Err()
	}

	sim, err := s.getSimulation(req.Id)
	if err != nil {
		return nil, status.Newf(codes.NotFound, "Failed to get simulation - %s", err).Err()
	}

	sim.Lock()
	defer sim.Unlock()

	if !sim.isRunning {
		return nil, status.Newf(codes.Unavailable, "Simulation Id: %s is not running", req.Id).Err()
	}

	sim.isRunning = false
	now := ServerTime()
	sim.stoppedTime = &now

	s.setSimulation(sim)

	DefaultLogger.log(fmt.Sprintf("Simulation: %s stop requested", sim.id))

	return resp, nil
}

func (s *server) NewSimulation(id, simName string, real *portfolio) (Simulation, error) {

	if _, ok := s.simulations[id]; ok {
		return nil, fmt.Errorf("Cannot create simulation %s as it already exists", id)
	}

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

func (s *simulation) GetID() string {
	return s.id
}

func (s *simulation) GetName() string {
	return s.name
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
	now := ServerTime()
	s.startedTime = &now
	s.Unlock()

	defer func() {
		s.Lock()
		s.isRunning = false
		s.Unlock()
	}()
	defer func() {
		s.Lock()
		now := ServerTime()
		s.stoppedTime = &now
		s.Unlock()
	}()

	if s.useHistoricalData {
		// TEMP: use some default dates for running simulation

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

	for priceTime := *s.simFromTime; priceTime.Before(toTime); priceTime = priceTime.Add(s.dataFrequency) {

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
	diff.print()

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
