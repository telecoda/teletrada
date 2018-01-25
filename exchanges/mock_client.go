package exchanges

import "time"

type mockClient struct {
}

const (
	MOCK_EXCHANGE = "mockexchange"
)

func NewMockClient() (ExchangeClient, error) {
	mock := &mockClient{}

	return mock, nil
}

func (m *mockClient) GetExchange() string {
	return MOCK_EXCHANGE
}

func (m *mockClient) GetBalances() ([]ExchangeBalance, error) {

	balances := []ExchangeBalance{
		ExchangeBalance{Symbol: "BTC", Free: 12.50, Locked: 12.50, Exchange: MOCK_EXCHANGE},
		ExchangeBalance{Symbol: "ETH", Free: 122.50, Locked: 122.50, Exchange: MOCK_EXCHANGE},
		ExchangeBalance{Symbol: "LTC", Free: 222.50, Locked: 222.50, Exchange: MOCK_EXCHANGE},
	}

	return balances, nil
}

func (m *mockClient) GetLatestPrices() ([]ExchangePrice, error) {

	mockPrices := []ExchangePrice{
		ExchangePrice{Base: "BTC", As: "USDT", Price: 12000.12345, At: time.Now()},
		ExchangePrice{Base: "ETH", As: "USDT", Price: 1.12345, At: time.Now()},
		ExchangePrice{Base: "LTC", As: "USDT", Price: 0.12345, At: time.Now()},
	}
	return mockPrices, nil
}

func (m *mockClient) GetHistoricPrices() ([]ExchangePrice, error) {
	mockOldPrices := []ExchangePrice{
		ExchangePrice{Base: "BTC", As: "USDT", Price: 13733.460000, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "BTC", As: "ETH", Price: 10000.12345, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "BTC", As: "BTC", Price: 9000.12345, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "ETH", As: "USDT", Price: 1.12346, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "ETH", As: "BTC", Price: 1.02345, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "ETH", As: "BTC", Price: 0.92345, At: time.Now().AddDate(0, 0, -1)},
		ExchangePrice{Base: "LTC", As: "ETH", Price: 0.12344, At: time.Now().AddDate(0, 0, -1)},
	}
	return mockOldPrices, nil
}
