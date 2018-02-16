package domain

import (
	"fmt"
	"os"

	"github.com/influxdata/influxdb/client/v2"
)

var DefaultMetrics *MetricsClient

const (
	INFLUX_DATABASE      = "teletrada"
	TEST_INFLUX_DATABASE = "testteletrada"
	INFLUX_HOST          = "http://localhost:8086"
)

type MetricsClient struct {
	client.Client
	dbName string
}

func newMetricsClient(dbName string) (*MetricsClient, error) {
	mc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     INFLUX_HOST,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating InfluxDB Client: %s", err.Error())
	}

	metricsClient := &MetricsClient{
		Client: mc,
		dbName: dbName,
	}

	return metricsClient, nil
}
