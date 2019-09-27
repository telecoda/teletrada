package domain

import (
	"fmt"
	"time"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

// This file contains the helper methods for the tests and mocking

var TEST_GRPC_PORT = 9999

// figures based on "rough" conversion rates on 5/5/2018

var _btcAsEth = 12.00     // x 12
var _btcAsBtc = 1.0       // 1:1
var _btcAsLtc = 55.00     // x 55
var _btcAsUsdt = 10000.00 // $10,000

var _ethAsBtc = 1.0 / _btcAsEth // divided by 12
var _ethAsEth = 1.0             // 1:1
var _ethAsLtc = 4.5             // x 4.5
var _ethAsUsdt = 800.00         // $800

var _ltcAsBtc = 1.0 / _btcAsLtc // divided by 55
var _ltcAsEth = 1.0 / _ethAsLtc
var _ltcAsLtc = 1.0     // 1:1
var _ltcAsUsdt = 180.00 // $180

func initMockServer() (Server, error) {

	servertime.InitFakeTime()
	// override server time func
	servertime.UseFakeTime()

	now := servertime.Now()

	// reset time
	// fakeTime gets bumped during setting up mock historical prices
	defer servertime.SetFakeTime(now)

	var err error

	config := Config{
		UseMock:        true,
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

// These are the mock balances used in all the simulation & portfolio tests. Be careful when you change them

func mockCoinBalances() []exchanges.CoinBalance {
	return []exchanges.CoinBalance{
		exchanges.CoinBalance{Symbol: "BTC", Free: 12.50, Locked: 12.50, Exchange: exchanges.MOCK_EXCHANGE},
		exchanges.CoinBalance{Symbol: "ETH", Free: 122.50, Locked: 122.50, Exchange: exchanges.MOCK_EXCHANGE},
		exchanges.CoinBalance{Symbol: "LTC", Free: 222.50, Locked: 222.50, Exchange: exchanges.MOCK_EXCHANGE},
	}
}

func goingUp(currentPrice, upwardsInc float64) func() float64 {
	// returns a funtion that starts at current price and will increase with each call
	return func() float64 {
		currentPrice += upwardsInc
		return currentPrice
	}
}

func goingDown(currentPrice, downwardsDec float64) func() float64 {
	// returns a funtion that starts at current price and will decrease with each call
	return func() float64 {
		currentPrice -= downwardsDec
		if currentPrice <= 0 {
			// default to min price
			currentPrice = 0.0000001
		}
		return currentPrice
	}
}

func goingUpDown(currentPrice, upDownInc, upDownMin, upDownMax float64, goingUp bool) func() float64 {
	return func() float64 {
		if goingUp {
			// going up
			currentPrice += upDownInc
			if currentPrice > upDownMax {
				currentPrice = upDownMax
				goingUp = false
			}
		} else {
			// going down
			currentPrice -= upDownInc
			if currentPrice < upDownMin {
				currentPrice = upDownMin
				goingUp = true
			}
		}
		return currentPrice
	}
}

// initMockPriceHistory - creates historic prices in archive and returns a slice of the latest prices
func initMockPriceHistory(when proto.StartSimulationRequestWhenOptions) ([]exchanges.Price, error) {

	toTime := servertime.Now()
	var fromTime time.Time

	switch when {
	case proto.StartSimulationRequest_LAST_DAY:
		fromTime = toTime.AddDate(0, 0, -1)
	case proto.StartSimulationRequest_LAST_WEEK:
		fromTime = toTime.AddDate(0, 0, -7)
	case proto.StartSimulationRequest_LAST_MONTH:
		fromTime = toTime.AddDate(0, 0, -30)
	case proto.StartSimulationRequest_THE_LOT:
		fromTime = toTime.AddDate(-10, 0, 0)
	default:
		return nil, fmt.Errorf("When value %d is not valid", when)
	}

	type baseAs struct {
		base  SymbolType
		as    SymbolType
		price func() float64
	}

	priceTypes := []baseAs{
		// BTC
		baseAs{base: BTC, as: BTC, price: func() float64 { return _btcAsBtc }},
		//baseAs{base: BTC, as: ETH, price: func() float64 { return _btcAsEth }},
		baseAs{base: BTC, as: ETH, price: goingUp(_btcAsEth, 0.00001)},
		baseAs{base: BTC, as: LTC, price: func() float64 { return _btcAsLtc }},
		baseAs{base: BTC, as: USDT, price: func() float64 { return _btcAsUsdt }},
		// ETH
		//baseAs{base: ETH, as: BTC, price: func() float64 { return _ethAsBtc }},
		baseAs{base: ETH, as: BTC, price: goingDown(_ethAsBtc, 0.00001)},
		baseAs{base: ETH, as: ETH, price: func() float64 { return _ethAsEth }},
		baseAs{base: ETH, as: LTC, price: func() float64 { return _ethAsLtc }},
		baseAs{base: ETH, as: USDT, price: func() float64 { return _ethAsUsdt }},
		// LTC
		baseAs{base: LTC, as: BTC, price: func() float64 { return _ltcAsBtc }},
		//baseAs{base: LTC, as: ETH, price: func() float64 { return _ltcAsEth }},
		baseAs{base: LTC, as: ETH, price: goingUpDown(_ltcAsEth, 0.01, _ltcAsEth/2.0, _ltcAsEth*2.0, true)},
		baseAs{base: LTC, as: LTC, price: func() float64 { return _ltcAsLtc }},
		baseAs{base: LTC, as: USDT, price: func() float64 { return _ltcAsUsdt }},
	}

	dataFrequency := time.Duration(1 * time.Minute)

	for priceTime := fromTime; priceTime.Before(toTime); priceTime = priceTime.Add(dataFrequency) {
		// update all symbols in portfolio

		for _, priceType := range priceTypes {
			price := Price{
				Base:     priceType.base,
				As:       priceType.as,
				At:       priceTime,
				Price:    priceType.price(),
				Exchange: "test-exchange",
			}
			err := DefaultArchive.AddPrice(price)
			if err != nil {
				return nil, err
			}
		}
	}

	// type priceFunc func() float64

	// 	upwardsPrice := 1.0
	// 	upwardsInc := 0.0001
	// 	downwardsPrice := 100.0
	// 	downwardsDec := 0.0001
	// 	upDownPrice := 50.0
	// 	upDownGoingUp := true
	// 	upDownMax := 75.00
	// 	upDownMin := 25.00
	// 	upDownInc := 0.02

	// 	goingUpAndDown := func() float64 {
	// 		if upDownGoingUp {
	// 			// going up
	// 			upDownPrice += upDownInc
	// 			if upDownPrice > upDownMax {
	// 				upDownPrice = upDownMax
	// 				upDownGoingUp = false
	// 			}
	// 		} else {
	// 			// going down
	// 			upDownPrice -= upDownInc
	// 			if upDownPrice < upDownMin {
	// 				upDownPrice = upDownMin
	// 				upDownGoingUp = true
	// 			}
	// 		}
	// 		return upDownPrice
	// 	}

	// 	pricers := map[SymbolType]priceFunc{
	// 		BTC: goingUp,
	// 		ETH: goingDown,
	// 		LTC: goingUpAndDown,
	// 	}

	// symbols := make([]SymbolType, len(s.livePortfolio.balances))
	// i := 0
	// for key, _ := range s.livePortfolio.balances {
	// 	symbols[i] = key
	// 	i++
	// }

	// accumulate latest prices from archive
	latestPrices := make([]exchanges.Price, 0)

	for _, priceType := range priceTypes {
		price, err := DefaultArchive.GetLatestPriceAs(priceType.base, priceType.as)
		if err != nil {
			return nil, err
		}

		latestPrices = append(latestPrices, price.toExchangePrice())
	}

	return latestPrices, nil
}

func initMockClients(prices []exchanges.Price) error {
	var err error
	DefaultClient, err = exchanges.NewMockClient(mockCoinBalances(), prices)
	if err != nil {
		return err
	}
	DefaultMetrics, err = newMockMetricsClient(TEST_INFLUX_DATABASE)
	if err != nil {
		return err
	}

	return nil
}
