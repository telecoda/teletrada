package domain

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
)

var TEST_GRPC_PORT = 9999

func setupTestServer() (Server, error) {
	var err error
	// Use mock metrics client to fetch price info
	DefaultMetrics, err = newMockMetricsClient(TEST_INFLUX_DATABASE)
	if err != nil {
		return nil, err
	}

	mc, err := exchanges.NewMockClient()
	if err != nil {
		return nil, err
	}

	// run test with mocked data
	DefaultClient = mc

	DefaultArchive = NewSymbolsArchive()

	config := Config{
		UseMock:        true,
		LoadPricesDir:  "",
		InfluxDBName:   "test-db",
		InfluxUsername: "",
		InfluxPassword: "",
		UpdateFreq:     time.Duration(1 * time.Hour),
		Verbose:        true,
		Port:           TEST_GRPC_PORT,
	}

	server, err := NewTradaServer(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create server - %s", err)
	}

	if err = server.Init(); err != nil {
		return nil, fmt.Errorf("Failed to init server - %s", err)
	}

	return server, nil
}
func TestStatusEndpoint(t *testing.T) {

	server, err := setupTestServer()
	assert.NoError(t, err)

	req := &proto.GetStatusRequest{}

	rsp, err := server.GetStatus(context.Background(), req)
	assert.NoError(t, err)
	assert.NotNil(t, rsp)

	assert.Equal(t, int32(3), rsp.TotalSymbols, "Test data should have 3 symbols")
	assert.NotZero(t, rsp.LastUpdate, "Should have updated on init")
	assert.NotZero(t, rsp.ServerStarted, "Should have a start time")
	assert.Equal(t, int32(1), rsp.UpdateCount, "Should have updated prices once at startup")

}
