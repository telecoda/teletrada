package exchanges

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/adshao/go-binance"
)

const (
	BINANCE_API_KEY    = "BINANCE_API_KEY"
	BINANCE_API_SECRET = "BINANCE_API_SECRET"
	BINANCE_EXCHANGE   = "binance"
)

const (
	BTC  = "BTC"
	BNB  = "BNB"
	ETH  = "ETH"
	LTC  = "LTC"
	USDT = "USDT"
)

type binanceClient struct {
	client *binance.Client
}

func NewBinanceClient() (ExchangeClient, error) {
	apiKey := os.Getenv(BINANCE_API_KEY)

	if apiKey == "" {
		return nil, fmt.Errorf("You must set environment variable %s with your key", BINANCE_API_KEY)
	}
	secretKey := os.Getenv(BINANCE_API_SECRET)
	if secretKey == "" {
		return nil, fmt.Errorf("You must set environment variable %s with your secret", BINANCE_API_SECRET)
	}

	client := &binanceClient{
		client: binance.NewClient(apiKey, secretKey),
	}

	return client, nil
}

func (b *binanceClient) GetExchange() string {
	return "binance"
}

func (b *binanceClient) GetCoinBalances() ([]CoinBalance, error) {
	res, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	balances := make([]CoinBalance, 0)
	for _, bal := range res.Balances {
		free, err := strconv.ParseFloat(bal.Free, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse balance quantity free: %s - %s. %#v", bal.Free, err, bal)
		}
		locked, err := strconv.ParseFloat(bal.Locked, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse balance quantity locked: %s - %s. %#v", bal.Locked, err, bal)
		}
		newBalance := CoinBalance{
			Symbol:   bal.Asset,
			Exchange: BINANCE_EXCHANGE,
			Free:     free,
			Locked:   locked,
		}
		if newBalance.Free != 0 || newBalance.Locked != 0 {
			balances = append(balances, newBalance)
		}
	}

	return balances, nil
}

func (b *binanceClient) GetLatestPrices() ([]Price, error) {
	res, err := b.client.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	prices := make([]Price, len(res))
	// convert results to prices
	for i, binancePrice := range res {
		symbolPrice, err := strconv.ParseFloat(binancePrice.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol price: %s - %s. %s", binancePrice.Symbol, err, binancePrice.Price)
		}

		// split LTCBTC symbol into base and as symbols
		// LTC/BTC

		price := Price{
			Price:    symbolPrice,
			At:       time.Now(),
			Exchange: BINANCE_EXCHANGE,
		}

		if strings.HasSuffix(binancePrice.Symbol, BTC) {
			// Bitcoin price
			price.As = BTC
		} else if strings.HasSuffix(binancePrice.Symbol, BNB) {
			// Binance coin price
			price.As = BNB
		} else if strings.HasSuffix(binancePrice.Symbol, ETH) {
			// Ether price
			price.As = ETH
		} else if strings.HasSuffix(binancePrice.Symbol, LTC) {
			// Litecoin price
			price.As = LTC
		} else if strings.HasSuffix(binancePrice.Symbol, USDT) {
			// US Dollar price
			price.As = USDT
		}

		if binancePrice.Symbol == "123456" {
			price.As = "123456"
			price.Base = "123456"
		} else {
			// only use price
			if price.As == "" {
				return nil, fmt.Errorf("Unexpected symbol type %s - %#v", binancePrice.Symbol, binancePrice)
			}
			price.Base = strings.Replace(binancePrice.Symbol, string(price.As), "", 1)

		}
		prices[i] = price

	}

	return prices, nil
}

func (b *binanceClient) GetHistoricPrices() ([]Price, error) {
	// TODO - get some old price data
	prices := make([]Price, 0)
	return prices, nil
}

func (b *binanceClient) GetDaySummaries() ([]DaySummary, error) {

	// Get latest prices for every coin
	info, err := b.client.NewExchangeInfoService().Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to get exchange info - %s", err)
	}

	// for each symbol, fetch last days price

	days := make([]DaySummary, len(info.Symbols))

	//errors := make([]error, 0)

	for i, symbol := range info.Symbols {
		key := symbol.BaseAsset + symbol.Symbol
		stats, err := b.client.NewPriceChangeStatsService().Symbol(key).Do(context.Background())
		if err != nil {
			// 	errors = append(errors, fmt.Errorf("Failed to get price change info for symbol %s - %s", symbol.BaseAsset, err))
			// 	continue
			return nil, fmt.Errorf("Failed to get price change info for symbol %s - %s", key, err)
		}

		days[i] = DaySummary{
			Base: symbol.BaseAsset,
			As:   symbol.Symbol,
		}

		openPrice, err := strconv.ParseFloat(stats.OpenPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol open price: %s - %s. %s", symbol.BaseAsset, stats.OpenPrice, err)
		}
		days[i].OpenPrice = openPrice

		closePrice, err := strconv.ParseFloat(stats.LastPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol close price: %s - %s. %s", symbol.BaseAsset, stats.LastPrice, err)
		}
		days[i].ClosePrice = closePrice

		weightedAvgPrice, err := strconv.ParseFloat(stats.WeightedAvgPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol weighted avg price: %s - %s. %s", symbol.BaseAsset, stats.WeightedAvgPrice, err)
		}
		days[i].WeightedAvgPrice = weightedAvgPrice

		highestPrice, err := strconv.ParseFloat(stats.HighPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol highest price: %s - %s. %s", symbol.BaseAsset, stats.HighPrice, err)
		}
		days[i].HighestPrice = highestPrice

		lowestPrice, err := strconv.ParseFloat(stats.LowPrice, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol lowest price: %s - %s. %s", symbol.BaseAsset, stats.LowPrice, err)
		}
		days[i].LowestPrice = lowestPrice

		changePrice, err := strconv.ParseFloat(stats.PriceChange, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol change price: %s - %s. %s", symbol.BaseAsset, stats.PriceChange, err)
		}
		days[i].ChangePrice = changePrice

		changePercent, err := strconv.ParseFloat(stats.PriceChangePercent, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol change percent: %s - %s. %s", symbol.BaseAsset, stats.PriceChangePercent, err)
		}
		days[i].ChangePercent = changePercent

		days[i].At = time.Unix(stats.CloseTime, 0)
		days[i].Exchange = b.GetExchange()
	}

	// // concat all error into one
	// if len(errors) > 0 {
	// 	errString := ""
	// 	for i, err := range errors {
	// 		errString += fmt.Sprintf("ERROR: (%d) - %s\n", i, err.Error())
	// 	}
	// 	return days, fmt.Errorf(errString)
	// }

	return days, nil
}

// func (b *binanceClient) GetPriceChange24(base, as string) (PriceChange24, error) {

// 	tp := base + as

// 	res, err := b.client.NewPriceChangeStatsService().Symbol(tp).Do(context.Background())
// 	if err != nil {
// 		return PriceChange24{}, err
// 	}

// 	priceChange := PriceChange24{
// 		Price: Price{
// 			Base:     base,
// 			As:       as,
// 			At:       time.Unix(0, res.CloseTime*1000000),
// 			Exchange: BINANCE_EXCHANGE,
// 		},
// 		OpenTime:  time.Unix(0, res.OpenTime*1000000),
// 		CloseTime: time.Unix(0, res.CloseTime*1000000),
// 	}

// 	lastPrice, err := strconv.ParseFloat(res.LastPrice, 64)
// 	if err != nil {
// 		return PriceChange24{}, fmt.Errorf("Failed to parse last price: %s - %s. %s", tp, err, res.LastPrice)
// 	}
// 	priceChange.Price.Price = lastPrice

// 	openPrice, err := strconv.ParseFloat(res.OpenPrice, 64)
// 	if err != nil {
// 		return PriceChange24{}, fmt.Errorf("Failed to parse open price: %s - %s. %s", tp, err, res.OpenPrice)
// 	}
// 	priceChange.OpenPrice = openPrice

// 	changeAmount, err := strconv.ParseFloat(res.PriceChange, 64)
// 	if err != nil {
// 		return PriceChange24{}, fmt.Errorf("Failed to parse price change: %s - %s. %s", tp, err, res.PriceChange)
// 	}
// 	priceChange.ChangeAmount = changeAmount

// 	changePercent, err := strconv.ParseFloat(res.PriceChangePercent, 64)
// 	if err != nil {
// 		return PriceChange24{}, fmt.Errorf("Failed to parse price change percent: %s - %s. %s", tp, err, res.PriceChange)
// 	}
// 	priceChange.ChangePercent = changePercent

// 	return priceChange, nil
// }
