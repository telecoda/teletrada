package domain

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/servertime"
	"google.golang.org/grpc/codes"
	sts "google.golang.org/grpc/status"
)

func createTestSimulation(s *server) (*simulation, error) {
	return s.newSimulation("test-sim-id", "test-sim-name")
}

func TestCreateSimulation(t *testing.T) {

	// test initialisation

	s, err := initMockServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// cast to internal type
	server := s.(*server)
	// end initialisation

	tests := []struct {
		name        string
		simName     string
		simID       string
		errExpected bool
		errText     string
	}{
		{
			name:        "Valid request",
			simName:     "Valid sim",
			simID:       "Valid sim ID",
			errExpected: false,
		},
		{
			name:        "Duplicate sim request",
			simName:     "Valid sim",
			simID:       "Valid sim ID",
			errExpected: true,
			errText:     "Cannot create simulation Valid sim ID as it already exists",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			req := &proto.CreateSimulationRequest{
				Id:   test.simID,
				Name: test.simName,
			}
			resp, err := server.CreateSimulation(ctx, req)

			if test.errExpected {
				assert.Error(t, err)
				assert.Nil(t, resp)
				if err != nil {
					assert.Equal(t, test.errText, err.Error())
				}
			} else {
				// check result
				if assert.NotNil(t, resp) {
					assert.NotNil(t, resp.Simulation)
					assert.Equal(t, test.simName, resp.Simulation.Name)
					assert.Equal(t, test.simID, resp.Simulation.Id)
				}
			}
		})
	}

}

func TestStartSimulationDates(t *testing.T) {

	// This test mainly checks that different simulation types are initialised with the correct dates
	servertime.UseFakeTime()
	defer servertime.UseRealTime()

	now := servertime.Now()

	tests := []struct {
		name              string
		when              proto.StartSimulationRequestWhenOptions
		simFromTime       time.Time
		simToTime         time.Time
		useHistoricalData bool
		useRealtimeData   bool
		errExpected       bool
		errText           string
	}{
		{
			name:              "Valid 1 day request",
			when:              proto.StartSimulationRequest_LAST_DAY,
			simFromTime:       now.AddDate(0, 0, -1),
			simToTime:         now,
			useHistoricalData: true,
			useRealtimeData:   false,
			errExpected:       false,
		},
		{
			name:              "Valid 1 week request",
			when:              proto.StartSimulationRequest_LAST_WEEK,
			simFromTime:       now.AddDate(0, 0, -7),
			simToTime:         now,
			useHistoricalData: true,
			useRealtimeData:   false,
			errExpected:       false,
		},
		{
			name:              "Valid 1 month request",
			when:              proto.StartSimulationRequest_LAST_MONTH,
			simFromTime:       now.AddDate(0, 0, -30),
			simToTime:         now,
			useHistoricalData: true,
			useRealtimeData:   false,
			errExpected:       false,
		},
		{
			name:              "Valid the lot request",
			when:              proto.StartSimulationRequest_THE_LOT,
			simFromTime:       now.AddDate(-10, 0, 0),
			simToTime:         now,
			useHistoricalData: true,
			useRealtimeData:   false,
			errExpected:       false,
		},
		{
			name:              "Valid realtime request",
			when:              proto.StartSimulationRequest_NOW_REALTIME,
			useHistoricalData: false,
			useRealtimeData:   true,
			errExpected:       false,
		},
		{
			name:              "Invalid request ",
			when:              proto.StartSimulationRequestWhenOptions(99999),
			useHistoricalData: false,
			useRealtimeData:   false,
			errExpected:       true,
			errText:           "When value 99999 is not valid",
		},
	}

	ctx := context.Background()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			s, err := initMockServer()
			assert.NoError(t, err)
			assert.NotNil(t, s)

			// cast to internal type
			server := s.(*server)

			// createTestSimulation
			sim, err := createTestSimulation(server)
			if assert.NoError(t, err) {
				assert.NotNil(t, sim)

				req := &proto.StartSimulationRequest{
					Id:   sim.id,
					When: test.when,
				}

				resp, err := server.StartSimulation(ctx, req)

				if test.errExpected {
					assert.Error(t, err)
					assert.Nil(t, resp)
					assert.Contains(t, err.Error(), test.errText)
				} else {
					assert.NoError(t, err)
					assert.NotNil(t, resp)
					// check times
					if sim.useHistoricalData {
						// only check sim dates on historical
						assert.Equal(t, test.simFromTime, *sim.simFromTime, "FromTime")
						assert.Equal(t, test.simToTime, *sim.simToTime, "ToTime")
					}
					assert.Equal(t, test.useHistoricalData, sim.useHistoricalData)
					assert.Equal(t, test.useRealtimeData, sim.useRealtimeData)
				}
			}
		})
	}
}

func TestSimulationStartStop(t *testing.T) {
	servertime.UseFakeTime()
	defer servertime.UseRealTime()

	s, err := initMockServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// Start a sim that doesn't exist

	startReq1 := &proto.StartSimulationRequest{
		Id:   "non-existent",
		When: proto.StartSimulationRequest_LAST_DAY,
	}

	ctx := context.Background()
	startResp1, err := s.StartSimulation(ctx, startReq1)

	assert.Nil(t, startResp1)
	assert.Error(t, err)
	status, ok := sts.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, status.Code())

	// Create a sim
	createReq := &proto.CreateSimulationRequest{
		Id:   "test-sim-001",
		Name: "test simulation number 1",
	}

	ctx = context.Background()
	createResp, err := s.CreateSimulation(ctx, createReq)
	assert.NotNil(t, createResp)
	assert.NoError(t, err)

	// Start a sim
	startReq2 := &proto.StartSimulationRequest{
		Id:   "test-sim-001",
		When: proto.StartSimulationRequest_LAST_DAY,
	}

	ctx = context.Background()
	startResp2, err := s.StartSimulation(ctx, startReq2)
	assert.NotNil(t, startResp2)
	assert.NoError(t, err)

	time.Sleep(100 * time.Millisecond)

	// Try to start it again
	startReq3 := &proto.StartSimulationRequest{
		Id:   "test-sim-001",
		When: proto.StartSimulationRequest_LAST_DAY,
	}

	ctx = context.Background()
	startResp3, err := s.StartSimulation(ctx, startReq3)
	assert.Nil(t, startResp3)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Simulation Id: test-sim-001 is already started")

	// Stop it
	stopReq1 := &proto.StopSimulationRequest{
		Id: "test-sim-001",
	}

	ctx = context.Background()
	stopResp1, err := s.StopSimulation(ctx, stopReq1)
	assert.NotNil(t, stopResp1)
	assert.NoError(t, err)

	// Try to stop it again
	stopReq2 := &proto.StopSimulationRequest{
		Id: "test-sim-001",
	}

	ctx = context.Background()
	stopResp2, err := s.StopSimulation(ctx, stopReq2)
	assert.Nil(t, stopResp2)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Simulation Id: test-sim-001 is not running")

}

func TestSimulationWithSimpleStrategy(t *testing.T) {

	servertime.UseFakeTime()
	defer servertime.UseRealTime()

	s, err := initMockServer()
	assert.NoError(t, err)
	assert.NotNil(t, s)

	// cast to internal type
	server := s.(*server)

	assert.NoError(t, err)

	// Create a sim
	createReq := &proto.CreateSimulationRequest{
		Id:   "test-sim-001",
		Name: "test simulation number 1",
	}

	ctx := context.Background()
	createResp, err := s.CreateSimulation(ctx, createReq)
	assert.NotNil(t, createResp)
	assert.NoError(t, err)

	testId := "test-sim-001"
	// Start a sim
	startReq := &proto.StartSimulationRequest{
		Id:   testId,
		When: proto.StartSimulationRequest_LAST_DAY,
	}

	ctx = context.Background()
	startResp, err := s.StartSimulation(ctx, startReq)
	assert.NotNil(t, startResp)
	assert.NoError(t, err)

	isFinished := false

	// keep checking until sim has finished running
	for !isFinished {

		req := &proto.GetSimulationsRequest{
			Id: testId,
		}

		resp, err := s.GetSimulations(ctx, req)
		assert.NoError(t, err)
		if err != nil {
			return
		}
		if len(resp.Simulations) != 1 {
			assert.Fail(t, "returned wrong number of simulations")
			return
		}
		sim := resp.Simulations[0]
		if sim.Id != testId {
			assert.Fail(t, "Id: %s does not match %s", testId, sim.Id)
		}

		if !sim.IsRunning && sim.StoppedTime != nil {
			isFinished = true

		} else {
			// still running
			// update fake time
			servertime.TickFakeTime(1 * time.Hour)
			// sleep for a bit
			time.Sleep(500 * time.Millisecond)
		}
	}

	// inspect results

	sim, ok := server.simulations[testId]
	if assert.True(t, ok, "Simulation %s not found", testId) {
		// only check results if actually found
		assert.Equal(t, testId, sim.id)
		assert.NotNil(t, sim.portfolio)
		assert.NotNil(t, sim.realAtStart)
		assert.NotNil(t, sim.realNow)

		assert.NotNil(t, sim.simFromTime)
		assert.NotNil(t, sim.simToTime)
		assert.NotNil(t, sim.startedTime)
		assert.NotNil(t, sim.stoppedTime)

		// 1 day test
		simDur := sim.simToTime.Sub(*sim.simFromTime)
		assert.Equal(t, time.Duration(24*time.Hour), simDur)
		execDur := sim.stoppedTime.Sub(*sim.startedTime)
		assert.Equal(t, time.Duration(1*time.Hour), execDur)

		// reprice "start" portfolio
		err := sim.realAtStart.repriceAt(*sim.simFromTime)
		assert.NoError(t, err)

		// reprice "now" portfolio
		err = sim.realNow.repriceAt(servertime.Now())
		assert.NoError(t, err)

		// reprice simulated portfolio
		//err = sim.portfolio.repriceAt(servertime.Now())
		//assert.NoError(t, err)

		// LTC as ETH goes up and down..

		// all balances are stored as BTC
		// need to convert prices to correct symbol
		// to do any real comparisons

		realEthStart, ok := sim.realAtStart.balances[ETH]
		if !ok {
			assert.Fail(t, "RealAtStart %s not found", ETH)
		}
		realEthNow, ok := sim.realNow.balances[ETH]
		if !ok {
			assert.Fail(t, "RealNow %s not found", ETH)
		}
		portEthNow, ok := sim.portfolio.balances[ETH]
		if !ok {
			assert.Fail(t, "Portfolio Now %s not found", ETH)
		}

		realBtcStart, ok := sim.realAtStart.balances[BTC]
		if !ok {
			assert.Fail(t, "RealAtStart %s not found", BTC)
		}
		realBtcNow, ok := sim.realNow.balances[BTC]
		if !ok {
			assert.Fail(t, "RealNow %s not found", BTC)
		}
		portBtcNow, ok := sim.portfolio.balances[BTC]
		if !ok {
			assert.Fail(t, "Portfolio Now %s not found", BTC)
		}

		// should have same number of coins still
		assert.Equal(t, realEthStart.Total, realEthNow.Total)
		assert.Equal(t, realBtcStart.Total, realBtcNow.Total)
		assert.Equal(t, realEthStart.Total, portEthNow.Total)
		assert.Equal(t, realBtcStart.Total, portBtcNow.Total)

		// should have different prices
		// ETH as BTC goes down
		if realEthStart.Price <= realEthNow.Price {
			fmt.Printf("ETH Start: %s - %#v\n", realEthStart.At.String(), realEthStart)
			fmt.Printf("ETH Now: %s - %#v\n", realEthNow.At.String(), realEthNow)
			assert.Fail(t, "ETH price didn't go down")
		}
		if realEthStart.Price <= portEthNow.Price {
			fmt.Printf("ServerTime: %s\n", servertime.Now().String())
			fmt.Printf("Sim From: %s\n", sim.simFromTime.String())
			fmt.Printf("Sim To: %s\n", sim.simToTime.String())
			fmt.Printf("Sim Started: %s\n", sim.startedTime.String())
			fmt.Printf("Sim Stopped: %s\n", sim.stoppedTime.String())
			fmt.Printf("ETH Start: %s - %f\n", realEthStart.At.String(), realEthStart.Price)
			fmt.Printf("ETH Now: %s - %f\n", realEthNow.At.String(), realEthNow.Price)
			fmt.Printf("ETH Port Now: %s - %f\n", portEthNow.At.String(), portEthNow.Price)
			assert.Fail(t, "ETH price didn't go down")
		}
		// should have different prices
		// BTC as ETH goes up

		// convert BTC as BTC to BTC as ETH
		realBtcStartAsEth, err := realBtcStart.convertTo(ETH)
		assert.NoError(t, err)

		realBtcNowAsEth, err := realBtcNow.convertTo(ETH)
		assert.NoError(t, err)

		portBtcNowAsEth, err := portBtcNow.convertTo(ETH)
		assert.NoError(t, err)

		if realBtcStartAsEth.Price >= realBtcNowAsEth.Price {
			fmt.Printf("BTC Start: %s - %#v\n", realBtcStartAsEth.At.String(), realBtcStartAsEth)
			fmt.Printf("BTC Now: %s - %#v\n", realBtcNowAsEth.At.String(), realBtcNowAsEth)
			assert.Fail(t, "BTC price didn't go up")
		}
		if realBtcStartAsEth.Price >= portBtcNowAsEth.Price {
			fmt.Printf("BTC Start: %s - %#v\n", realBtcStartAsEth.At.String(), realBtcStartAsEth)
			fmt.Printf("BTC Now: %s - %#v\n", realBtcNowAsEth.At.String(), realBtcNowAsEth)
			fmt.Printf("BTC Port Now: %s - %#v\n", portBtcNowAsEth.At.String(), portBtcNowAsEth)
			assert.Fail(t, "BTC price didn't go up")
		}
		// check if strategy has ever triggered

	}

}
