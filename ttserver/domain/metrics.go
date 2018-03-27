package domain

import (
	"fmt"
	"log"
	"os"

	"github.com/influxdata/influxdb/client/v2"
)

var DefaultMetrics MetricsClient

const (
	INFLUX_DATABASE      = "teletrada"
	TEST_INFLUX_DATABASE = "testteletrada"
	INFLUX_HOST          = "http://localhost:8086"
)

type MetricsClient interface {
	GetDBName() string
	SavePriceMetrics(prices []Price) error
	SavePortfolioMetrics(portfolio *portfolio) error
}

type metricsClient struct {
	client.Client
	dbName string
}

func newMetricsClient(dbName string) (MetricsClient, error) {
	mc, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     INFLUX_HOST,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	})
	if err != nil {
		return nil, fmt.Errorf("Error creating InfluxDB Client: %s", err.Error())
	}

	mClient := &metricsClient{
		Client: mc,
		dbName: dbName,
	}

	return mClient, nil
}

func (m *metricsClient) GetDBName() string {
	return m.dbName
}

func (m *metricsClient) SavePriceMetrics(prices []Price) error {

	if len(prices) == 0 {
		return nil
	}

	log.Printf("Sending symbol price data to influxdb")
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  m.dbName,
		Precision: "ns",
	})

	if err != nil {
		return fmt.Errorf("failed to create batch points: %s", err)
	}
	for _, price := range prices {

		if price.As == "123456" {
			continue // skip it
		}

		// Create a point and add to batch
		tags := map[string]string{"symbol": string(price.Base)}
		fields := make(map[string]interface{}, 0)

		toSymbols := []SymbolType{SymbolType(BTC), SymbolType(ETH), SymbolType(USDT)}

		for _, toSym := range toSymbols {
			if symPrice, err := DefaultArchive.GetLatestPriceAs(price.Base, toSym); err != nil {
				log.Printf("No %s price for %s symbol - %s", toSym, price.Base, err)
			} else {
				fields[fmt.Sprintf("price.%s", toSym)] = symPrice.Price
			}
		}

		if len(fields) > 0 {
			fields["exchange"] = price.Exchange
			// only add fields with points
			pt, err := client.NewPoint("coin_price", tags, fields, price.At)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}

			bp.AddPoint(pt)
		}
	}
	// Write the batch
	if err := m.Write(bp); err != nil {
		log.Printf("error sending metrics %s", err)
		return err
	}

	return nil
}

func (m *metricsClient) SavePortfolioMetrics(p *portfolio) error {

	log.Printf("Sending portfolio balance data to influxdb")
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  m.dbName,
		Precision: "ns",
	})

	if err != nil {
		return fmt.Errorf("failed to create batch points: %s", err)
	}

	portType := "live"
	if !p.isLive {
		portType = "simulated"
	}

	for _, balance := range p.balances {

		// Create a point and add to batch
		tags := map[string]string{"symbol": string(balance.Symbol), "name": p.name,
			"live": portType}
		fields := make(map[string]interface{}, 0)

		toSymbols := []SymbolType{SymbolType(BTC), SymbolType(ETH), SymbolType(USDT)}

		for _, toSym := range toSymbols {
			if symPrice, err := DefaultArchive.GetLatestPriceAs(SymbolType(balance.Symbol), toSym); err != nil {
				log.Printf("No %s price for %s symbol - %s", toSym, balance.Symbol, err)
			} else {
				fields[fmt.Sprintf("price.%s", toSym)] = symPrice.Price
				// calc current value = total * price
				value := symPrice.Price * balance.Total
				fields[fmt.Sprintf("value.%s", toSym)] = value
			}
		}

		if len(fields) > 0 {
			fields["exchange"] = balance.Exchange
			// add coin totals
			fields["total"] = balance.Total
			fields["locked"] = balance.Locked
			fields["free"] = balance.Free

			// only add fields with points
			pt, err := client.NewPoint("portfolio_balance", tags, fields, balance.At)
			if err != nil {
				log.Printf("Error: %s", err.Error())
			} else {
				bp.AddPoint(pt)
			}

		}
	}
	// Write the batch
	return m.Write(bp)

}

type mockMetricsClient struct {
	dbName string
}

func newMockMetricsClient(dbName string) (MetricsClient, error) {

	mClient := &mockMetricsClient{
		dbName: dbName,
	}

	return mClient, nil
}

func (m *mockMetricsClient) GetDBName() string {
	return m.dbName
}

func (m *mockMetricsClient) SavePriceMetrics(prices []Price) error {
	return nil
}
func (m *mockMetricsClient) SavePortfolioMetrics(portfolio *portfolio) error {
	return nil
}
