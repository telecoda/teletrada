package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/proto"

	"github.com/telecoda/teletrada/exchanges"
)

func createTestPortfolio() *portfolio {

	p := &portfolio{
		name:     "Live",
		isLive:   true,
		balances: make(map[SymbolType]*BalanceAs),
	}

	now := ServerTime()
	// add some balances
	balance1 := &BalanceAs{
		CoinBalance: exchanges.CoinBalance{
			Symbol:   "symbol1",
			Exchange: "test-exchange",
			Free:     5,
			Locked:   5,
		},
		Total:        10,
		At:           now,
		As:           SymbolType("BTC"),
		Price:        25.00,
		Value:        250.00,
		Price24H:     20.00,
		Value24H:     200.00,
		Change24H:    50.00,
		ChangePct24H: 25.0,
	}
	balance2 := &BalanceAs{
		CoinBalance: exchanges.CoinBalance{
			Symbol:   "symbol2",
			Exchange: "test-exchange",
			Free:     50,
			Locked:   50,
		},
		Total:        10,
		At:           now,
		As:           SymbolType("BTC"),
		Price:        25.00,
		Value:        2500.00,
		Price24H:     20.00,
		Value24H:     2000.00,
		Change24H:    500.00,
		ChangePct24H: 25.0,
	}

	p.balances["symbol1"] = balance1
	p.balances["symbol2"] = balance2

	return p

}

func createTestSimulation(s *server) (Simulation, error) {
	realPort := createTestPortfolio()
	return s.NewSimulation("test-sim-id", "test-sim-name", realPort)
}

func TestNewSimulation(t *testing.T) {

	// test initialisation

	s, err := setupTestServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	real := createTestPortfolio()
	// cast to internal type
	server := s.(*server)

	// end initialisation

	tests := []struct {
		name        string
		simName     string
		simID       string
		portfolio   *portfolio
		errExpected bool
		errText     string
	}{
		{
			name:        "Valid request",
			simName:     "Valid sim",
			simID:       "Valid sim ID",
			portfolio:   real,
			errExpected: false,
		},
		{
			name:        "Duplicate sim request",
			simName:     "Valid sim",
			simID:       "Valid sim ID",
			portfolio:   real,
			errExpected: true,
			errText:     "Cannot create simulation Valid sim ID as it already exists",
		},
		{
			name:        "Nil portfolio request",
			simName:     "nil port sim",
			simID:       "nil port sim ID",
			portfolio:   nil,
			errExpected: true,
			errText:     "Portfolio cannot be nil",
		},
		{
			name:        "Empty portfolio request",
			simName:     "nil port sim",
			simID:       "nil port sim ID",
			portfolio:   &portfolio{},
			errExpected: true,
			errText:     "Cannot use a portfolio with no balances for a simulation",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sim, err := server.NewSimulation(test.simID, test.simName, test.portfolio)

			if test.errExpected {
				assert.Error(t, err)
				assert.Nil(t, sim)
				assert.Equal(t, test.errText, err.Error())
			} else {
				// check result
				assert.NotNil(t, sim)
				assert.Equal(t, test.simName, sim.GetName())
				assert.Equal(t, test.simID, sim.GetID())
			}
		})
	}

}

func TestStartSimulation(t *testing.T) {

	UseFakeTime()
	defer UseRealTime()

	now := ServerTime()

	tests := []struct {
		name        string
		when        proto.StartSimulationRequestWhenOptions
		simFromTime time.Time
		simToTime   time.Time
		errExpected bool
		errText     string
	}{
		{
			name:        "Valid 1 day request",
			when:        proto.StartSimulationRequest_LAST_DAY,
			simFromTime: now.AddDate(0, 0, -1),
			simToTime:   now,
			errExpected: false,
		},
		{
			name:        "Valid 1 week request",
			when:        proto.StartSimulationRequest_LAST_WEEK,
			simFromTime: now.AddDate(0, 0, -7),
			simToTime:   now,
			errExpected: false,
		},
		{
			name:        "Valid 1 month request",
			when:        proto.StartSimulationRequest_LAST_MONTH,
			simFromTime: now.AddDate(0, 0, -30),
			simToTime:   now,
			errExpected: false,
		},
		{
			name:        "Valid the lot request",
			when:        proto.StartSimulationRequest_THE_LOT,
			simFromTime: now.AddDate(-10, 0, 0),
			simToTime:   now,
			errExpected: false,
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s, err := setupTestServer()
			assert.NoError(t, err)
			assert.NotNil(t, s)

			// cast to internal type
			server := s.(*server)

			// createTestSimulation
			sim, err := createTestSimulation(server)
			if assert.NoError(t, err) {
				assert.NotNil(t, sim)

				req := &proto.StartSimulationRequest{
					Id:   sim.GetID(),
					When: test.when,
				}

				resp, err := server.StartSimulation(ctx, req)
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				// check times
				// cast to native type
				realSim := sim.(*simulation)
				assert.Equal(t, test.simFromTime, *realSim.simFromTime, "FromTime")
				assert.Equal(t, test.simToTime, *realSim.simToTime, "ToTime")
			}
		})
	}
}
