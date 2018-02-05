package domain

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

type SymbolType string

type Symbol interface {
	AddPrice(price Price)
	GetType() SymbolType
	GetPriceAs(symbol SymbolType, at time.Time) (Price, error)
	GetLatestPriceAs(symbol SymbolType) (Price, error)
}

type symbol struct {
	sync.RWMutex
	SymbolType
	// map of prices by currency
	// etc for LTC symbol it may have prices for
	// LTCBTC, LTCETH and LTCUSDT
	priceAs map[SymbolType][]Price
}

func NewSymbol(symbolType SymbolType) *symbol {
	return &symbol{
		SymbolType: symbolType,
		priceAs:    make(map[SymbolType][]Price),
	}
}

func (s *symbol) GetType() SymbolType {
	return s.SymbolType
}

func (s *symbol) AddPrice(price Price) {
	s.Lock()
	defer s.Unlock()
	var prices []Price
	prices, _ = s.priceAs[price.As]
	prices = append(prices, price)

	// sort in date order
	sort.Slice(prices, func(i, j int) bool { return prices[i].At.Before(prices[j].At) })

	s.priceAs[price.As] = prices
}

// GetPriceAs - returns the price of base symbol as another symbol at a particular time
func (s *symbol) GetPriceAs(as SymbolType, at time.Time) (Price, error) {
	s.RLock()
	defer s.RUnlock()

	// base and as symbol are the same so price == 1
	if s.SymbolType == as {
		return Price{
			Base:  s.SymbolType,
			As:    as,
			Price: 1.0,
			At:    at,
		}, nil
	}

	if prices, ok := s.priceAs[as]; !ok || len(prices) == 0 {
		return Price{}, fmt.Errorf("Symbol: %s has no price information for: %s", s.SymbolType, as)
	} else {
		// return search for price with nearest date
		var priceBefore Price
		var priceAfter Price

		// read backwards through prices
		for i := len(prices) - 1; i > -1; i-- {

			price := prices[i]
			// if price is still later than requested time
			// save details
			if price.At.After(at) || price.At.Equal(at) {
				priceAfter = price
				priceBefore = price
			} else {
				// price must be before request time
				priceBefore = price
				break
			}

		}

		/*

			How time calcs work:

				 price 1  01:00:00 £1000.00
				 price 2  05:00:00 £50000.00

				request time @ 1:00:00 get price £1000.00
				request time @ 5:00:00 get price £5000.00
				request time @ 2:00:00 get price £2000.00

				eg. checking a 2pm price

				betweenPrices = afterDate - beforeDate

				05:00:00 - 01:00:00 = 4 hours

				sinceBefore = priceAt - beforeDate

				02:00:00 - 01:00:00 = 1 hour

				priceChange = afterPrice - beforePrice

				5000.00 - 1000.00 = 4000.00

				ratio = sinceBefore / betweenPrices

				1 hour / 4 hours = 1/4

				adjustedPrice = beforePrice + (priceChange * ratio)

				1000.00 + (1/4 * 4000.00) = 2000.00


		*/

		betweenPrices := priceAfter.At.Sub(priceBefore.At)

		sinceBefore := at.Sub(priceBefore.At)

		// before & after are the same
		if betweenPrices == 0 {
			priceAfter.At = at
			return priceAfter, nil
		}

		// adjust price to be between the two prices
		priceChange := priceAfter.Price - priceBefore.Price

		ratio := float64(sinceBefore.Nanoseconds()) / float64(betweenPrices.Nanoseconds())

		adjustedPrice := priceBefore.Price + (priceChange * ratio)

		priceAdjusted := priceBefore
		priceAdjusted.Price = adjustedPrice
		priceAdjusted.At = at
		return priceAdjusted, nil

	}
}

// GetLatestPriceAs - returns the latest price of base symbol as another symbol
func (s *symbol) GetLatestPriceAs(as SymbolType) (Price, error) {
	s.RLock()
	defer s.RUnlock()

	// base and as symbol are the same so price == 1
	if s.SymbolType == as {
		return Price{
			Base:  s.SymbolType,
			As:    as,
			Price: 1.0,
			At:    time.Now().UTC(),
		}, nil
	}
	if prices, ok := s.priceAs[as]; !ok || len(prices) == 0 {
		return Price{}, fmt.Errorf("Symbol: %s has no price information for: %s", s.SymbolType, as)
	} else {
		// return last price
		return prices[len(prices)-1], nil
	}
}
