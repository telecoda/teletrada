package domain

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

type SymbolType string

type Symbol interface {
	GetType() SymbolType
	GetAsTypes() []SymbolType
	// Prices
	AddPrice(price Price)
	GetPriceAs(as SymbolType, at time.Time) (Price, error)
	GetLatestPriceAs(as SymbolType) (Price, error)
	// Daily summary
	AddDaySummary(sum DaySummary)
	GetDaySummaryAs(as SymbolType) (DaySummary, error)
}

type symbol struct {
	sync.RWMutex
	SymbolType
	// map of prices by currency
	// etc for LTC symbol it may have prices for
	// LTCBTC, LTCETH and LTCUSDT
	priceAs    map[SymbolType][]Price // an array of prices by symbol type
	daySummary map[SymbolType]DaySummary
}

func NewSymbol(symbolType SymbolType) *symbol {
	return &symbol{
		SymbolType: symbolType,
		priceAs:    make(map[SymbolType][]Price),
		daySummary: make(map[SymbolType]DaySummary),
	}
}

func (s *symbol) GetType() SymbolType {
	s.RLock()
	defer s.RUnlock()
	return s.SymbolType
}

func (s *symbol) GetAsTypes() []SymbolType {
	s.RLock()
	defer s.RUnlock()
	asTypes := make([]SymbolType, len(s.priceAs))
	i := 0
	for k, _ := range s.priceAs {
		asTypes[i] = k
		i++
	}
	return asTypes
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

func (s *symbol) AddDaySummary(sum DaySummary) {
	s.Lock()
	defer s.Unlock()
	s.daySummary[sum.As] = sum
}

func (s *symbol) GetDaySummaryAs(as SymbolType) (DaySummary, error) {
	s.RLock()
	defer s.RUnlock()
	if sum, ok := s.daySummary[as]; !ok {
		return DaySummary{}, fmt.Errorf("Symbol: %s has no daily summary for: %s", s.SymbolType, as)
	} else {
		return sum, nil
	}
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
			At:    servertime.Now(),
		}, nil
	}
	if prices, ok := s.priceAs[as]; !ok || len(prices) == 0 {
		return Price{}, fmt.Errorf("Symbol: %s has no price information for: %s", s.SymbolType, as)
	} else {
		// return last price
		return prices[len(prices)-1], nil
	}
}

// GetSymbolTypes returns list of available symbols
func (s *server) GetSymbolTypes(ctx context.Context, req *proto.GetSymbolTypesRequest) (*proto.GetSymbolTypesResponse, error) {

	symbols := DefaultArchive.GetSymbolTypes()

	resp := &proto.GetSymbolTypesResponse{
		SymbolTypes: make([]*proto.SymbolType, len(symbols)),
	}

	i := 0
	for k, asTypes := range symbols {
		protoSymbol := &proto.SymbolType{}
		protoSymbol.Base = string(k)
		protoSymbol.As = make([]string, len(asTypes))
		for n, asType := range asTypes {
			protoSymbol.As[n] = string(asType)
		}
		resp.SymbolTypes[i] = protoSymbol
		i++
	}

	return resp, nil
}
