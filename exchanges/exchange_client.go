package exchanges

import "time"

type ExchangeClient interface {
	GetBalances() ([]ExchangeBalance, error)
	GetLatestPrices() ([]ExchangePrice, error)
	GetHistoricPrices() ([]ExchangePrice, error)
	GetExchange() string
}

type ExchangeBalance struct {
	Symbol   string
	Exchange string
	Free     float64
	Locked   float64
}

type ExchangePrice struct {
	Base  string
	As    string
	Price float64
	At    time.Time
}
