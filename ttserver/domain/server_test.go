package domain

import (
	"fmt"
	"time"

	"github.com/telecoda/teletrada/exchanges"
)

// This file contains the helper methods for the tests

var TEST_GRPC_PORT = 9999

func setupTestServer() (Server, error) {

	// override server time func
	ServerTime = fakeServerTime

	var err error
	// Use mock metrics client to fetch price info
	DefaultMetrics, err = newMockMetricsClient(TEST_INFLUX_DATABASE)
	if err != nil {
		return nil, err
	}

	// Use mock exchange client - this will return mock balances
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

	// Init the server will initialise the dummy portfolio
	if err = server.Init(); err != nil {
		return nil, fmt.Errorf("Failed to init server - %s", err)
	}

	return server, nil
}
