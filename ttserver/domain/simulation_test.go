package domain

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/proto"
	"google.golang.org/grpc/codes"
	sts "google.golang.org/grpc/status"
)

func createTestSimulation(s *server) (*simulation, error) {
	return s.newSimulation("test-sim-id", "test-sim-name")
}

func TestCreateSimulation(t *testing.T) {

	// test initialisation

	s, err := setupTestServer()
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
				assert.NotNil(t, resp)
				assert.Equal(t, test.simName, resp.Simulation.Name)
				assert.Equal(t, test.simID, resp.Simulation.Id)
			}
		})
	}

}

func TestStartSimulationDates(t *testing.T) {

	// This test mainly checks that different simulation types are initialised with the correct dates
	UseFakeTime()
	defer UseRealTime()

	now := ServerTime()

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
	UseFakeTime()
	defer UseRealTime()

	s, err := setupTestServer()
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
