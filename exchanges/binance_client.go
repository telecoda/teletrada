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
		return nil, fmt.Errorf("You must set environment variable %s with you key", BINANCE_API_KEY)
	}
	secretKey := os.Getenv(BINANCE_API_SECRET)
	if secretKey == "" {
		return nil, fmt.Errorf("You must set environment variable %s with you key", BINANCE_API_SECRET)
	}

	client := &binanceClient{
		client: binance.NewClient(apiKey, secretKey),
	}

	return client, nil
}

func (b *binanceClient) GetExchange() string {
	return "binance"
}

func (b *binanceClient) GetBalances() ([]ExchangeBalance, error) {
	res, err := b.client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	balances := make([]ExchangeBalance, 0)
	for _, bal := range res.Balances {
		free, err := strconv.ParseFloat(bal.Free, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse balance quantity free: %s - %s. %#v", bal.Free, err, bal)
		}
		locked, err := strconv.ParseFloat(bal.Locked, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse balance quantity locked: %s - %s. %#v", bal.Locked, err, bal)
		}
		newBalance := ExchangeBalance{
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

func (b *binanceClient) GetLatestPrices() ([]ExchangePrice, error) {
	res, err := b.client.NewListPricesService().Do(context.Background())
	if err != nil {
		return nil, err
	}

	prices := make([]ExchangePrice, len(res))
	// convert results to prices
	for i, binancePrice := range res {
		symbolPrice, err := strconv.ParseFloat(binancePrice.Price, 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse symbol price: %s - %s. %s", binancePrice.Symbol, err, binancePrice.Price)
		}

		// split LTCBTC symbol into base and as symbols
		// LTC/BTC

		price := ExchangePrice{
			Price: symbolPrice,
			At:    time.Now(),
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

func (b *binanceClient) GetHistoricPrices() ([]ExchangePrice, error) {
	// TODO - get some old price data
	prices := make([]ExchangePrice, 0)
	return prices, nil
}
