package exchanges

import (
	"time"
)

type mockClient struct {
	balances     []CoinBalance
	prices       []Price
	daySummaries []DaySummary
	serverTime   func() time.Time
}

const (
	MOCK_EXCHANGE = "mockexchange"
)

func NewMockClient(coinBalances []CoinBalance, exchangePrices []Price) (ExchangeClient, error) {
	mock := &mockClient{
		balances: coinBalances,
		prices:   exchangePrices,
	}

	return mock, nil
}

func (m *mockClient) GetExchange() string {
	return MOCK_EXCHANGE
}

// func initMockDaySummaries() []DaySummary {
// 	return []DaySummary{
// 		DaySummary{
// 			Base:             "BTC",
// 			As:               "USDT",
// 			OpenPrice:        10000.00,
// 			ClosePrice:       11000.00,
// 			WeightedAvgPrice: 10500.00,
// 			HighestPrice:     11500.00,
// 			LowestPrice:      9500.00,
// 			ChangePrice:      1000.00,
// 			ChangePercent:    10.00,
// 			At:               servertime.Now(),
// 			Exchange:         MOCK_EXCHANGE,
// 		},
// 		DaySummary{
// 			Base:             "ETH",
// 			As:               "USDT",
// 			OpenPrice:        1100.00,
// 			ClosePrice:       1000.00,
// 			WeightedAvgPrice: 1050.00,
// 			HighestPrice:     1150.00,
// 			LowestPrice:      950.00,
// 			ChangePrice:      -100.00,
// 			ChangePercent:    -10.00,
// 			At:               servertime.Now(),
// 			Exchange:         MOCK_EXCHANGE,
// 		},
// 		DaySummary{
// 			Base:             "LTC",
// 			As:               "USDT",
// 			OpenPrice:        110.00,
// 			ClosePrice:       100.00,
// 			WeightedAvgPrice: 150.00,
// 			HighestPrice:     150.00,
// 			LowestPrice:      90.00,
// 			ChangePrice:      -10.00,
// 			ChangePercent:    -99.1234,
// 			At:               servertime.Now(),
// 			Exchange:         MOCK_EXCHANGE,
// 		},
// 	}
// }

func (m *mockClient) GetCoinBalances() ([]CoinBalance, error) {
	return m.balances, nil
}

func (m *mockClient) GetLatestPrices() ([]Price, error) {
	return m.prices, nil
}

// func (m *mockClient) GetHistoricPrices() ([]Price, error) {
// 	mockOldPrices := []Price{
// 		Price{Base: "BTC", As: "USDT", Price: 13733.460000, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "BTC", As: "ETH", Price: 10000.12345, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "BTC", As: "BTC", Price: 9000.12345, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "ETH", As: "USDT", Price: 1.12346, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "ETH", As: "BTC", Price: 1.02345, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "ETH", As: "BTC", Price: 0.92345, At: time.Now().AddDate(0, 0, -1)},
// 		Price{Base: "LTC", As: "ETH", Price: 0.12344, At: time.Now().AddDate(0, 0, -1)},
// 	}
// 	return mockOldPrices, nil
// }

func (m *mockClient) GetDaySummaries() ([]DaySummary, error) {
	return m.daySummaries, nil
}
