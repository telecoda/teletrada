package exchanges

import "time"

type ExchangeClient interface {
	GetBalances() ([]Balance, error)
	GetLatestPrices() ([]Price, error)
	GetHistoricPrices() ([]Price, error)
	GetDaySummaries() ([]DaySummary, error)
	GetExchange() string
}

type Balance struct {
	Symbol   string
	Exchange string
	Free     float64
	Locked   float64
}

type Price struct {
	Base     string // This is the base symbol eg. NEO
	As       string // This is trading pair symbol eg. BTC, ETH etc
	Price    float64
	Exchange string
	At       time.Time
}

type DaySummary struct {
	Base             string
	As               string
	OpenPrice        float64
	ClosePrice       float64
	WeightedAvgPrice float64
	HighestPrice     float64
	LowestPrice      float64
	ChangePrice      float64
	ChangePercent    float64
	At               time.Time
	Exchange         string
}

/*
{
	"symbol": "BTCUSDT",
	"priceChange": "405.35000000",
	"priceChangePercent": "3.534",
	"weightedAvgPrice": "11855.42965188",
	"prevClosePrice": "11497.99000000",
	"lastPrice": "11876.06000000",
	"lastQty": "0.08374300",
	"bidPrice": "11875.50000000",
	"bidQty": "0.00088100",
	"askPrice": "11875.70000000",
	"askQty": "0.00010100",
	"openPrice": "11470.71000000",
	"highPrice": "12244.00000000",
	"lowPrice": "11408.00000000",
	"volume": "16891.27751600",
	"quoteVolume": "200253352.32137449",
	"openTime": 1517097500866,
	"closeTime": 1517183900866,
	"firstId": 10640313,
	"lastId": 10846665,
	"count": 206353
},
*/

// type PriceChange24 struct {
// 	Price
// 	ChangePercent float64
// 	ChangeAmount  float64
// 	OpenPrice     float64
// 	OpenTime      time.Time
// 	CloseTime     time.Time
// }
