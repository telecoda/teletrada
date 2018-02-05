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

func (m *mockClient) GetBalances() ([]Balance, error) {

	balances := []Balance{
		Balance{Symbol: "BTC", Free: 12.50, Locked: 12.50, Exchange: MOCK_EXCHANGE},
		Balance{Symbol: "ETH", Free: 122.50, Locked: 122.50, Exchange: MOCK_EXCHANGE},
		Balance{Symbol: "LTC", Free: 222.50, Locked: 222.50, Exchange: MOCK_EXCHANGE},
	}

	return balances, nil
}

func (m *mockClient) GetLatestPrices() ([]Price, error) {

	mockPrices := []Price{
		Price{Base: "BTC", As: "USDT", Price: 12000.12345, At: time.Now()},
		Price{Base: "LTC", As: "ETH", Price: 0.12345, At: time.Now()},
	}
	return mockPrices, nil
}

func (m *mockClient) GetHistoricPrices() ([]Price, error) {
	mockOldPrices := []Price{
		Price{Base: "BTC", As: "USDT", Price: 13733.460000, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "BTC", As: "ETH", Price: 10000.12345, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "BTC", As: "BTC", Price: 9000.12345, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "ETH", As: "USDT", Price: 1.12346, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "ETH", As: "BTC", Price: 1.02345, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "ETH", As: "BTC", Price: 0.92345, At: time.Now().AddDate(0, 0, -1)},
		Price{Base: "LTC", As: "ETH", Price: 0.12344, At: time.Now().AddDate(0, 0, -1)},
	}
	return mockOldPrices, nil
}

func (m *mockClient) GetPriceChange24(base, as string) (PriceChange24, error) {
	now := time.Now().UTC()
	mockPriceChange := PriceChange24{
		Price: Price{
			Base:  base,
			As:    as,
			Price: 1250.00,
			At:    now,
		},
		ChangePercent: 25.00,
		ChangeAmount:  250.00,
		OpenPrice:     1000.00,
		OpenTime:      now.AddDate(0, 0, -1),
		CloseTime:     now,
	}
	return mockPriceChange, nil
}
