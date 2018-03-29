package domain

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/telecoda/teletrada/proto"
)

type Price struct {
	Base     SymbolType
	As       SymbolType
	Price    float64
	At       time.Time
	Exchange string
}

type DaySummary struct {
	Base             SymbolType
	As               SymbolType
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

func (p Price) Validate() error {
	if p.Base == "" {
		return fmt.Errorf("Price invalid: Base symbol cannot be blank")
	}
	if p.As == "" {
		return fmt.Errorf("Price invalid: As symbol cannot be blank")
	}
	if p.Price == 0 {
		return fmt.Errorf("Price invalid: Price cannot be zero")
	}
	if p.At.IsZero() {
		return fmt.Errorf("Price invalid: At cannot be zero")
	}

	return nil

}

// GetPrices returns current prices
func (s *server) GetPrices(ctx context.Context, req *proto.GetPricesRequest) (*proto.GetPricesResponse, error) {
	resp := &proto.GetPricesResponse{}

	req.Base = strings.ToUpper(req.Base)
	req.As = strings.ToUpper(req.As)

	var symbolTypes []SymbolType

	if req.Base == "" || req.Base == "*ALL" {
		// all prices
		symbolMap := DefaultArchive.GetSymbolTypes()
		symbolTypes = make([]SymbolType, len(symbolMap))
		i := 0
		for symbolType, _ := range symbolMap {
			symbolTypes[i] = symbolType
			i++
		}

	} else {
		// only one symbol
		symbolTypes = make([]SymbolType, 1)
		symbolTypes[0] = SymbolType(req.Base)
	}

	resp.Prices = make([]*proto.Price, len(symbolTypes))

	for i, symbolType := range symbolTypes {
		price, err := DefaultArchive.GetLatestPriceAs(symbolType, SymbolType(req.As))
		if err != nil {
			return nil, fmt.Errorf("Failed to fetch symbol %s price as %s - %s", req.Base, req.As, err)
		}

		pp, err := price.toProto()
		if err != nil {
			return nil, err
		}

		daySummary, err := DefaultArchive.GetDaySummaryAs(symbolType, SymbolType(req.As))
		if err == nil {
			// day summary found so fill in corresponding fields
			pp.ChangePct24H = float32(daySummary.ChangePercent)
			pp.Change24H = float32(daySummary.ChangePrice)
			pp.Opening = float32(daySummary.OpenPrice)
			pp.Closing = float32(daySummary.ClosePrice)
			pp.Highest = float32(daySummary.HighestPrice)
			pp.Lowest = float32(daySummary.LowestPrice)
			pp.ChangeToday = pp.Current - pp.Closing
			if pp.ChangeToday != 0 {
				pp.ChangePctToday = (pp.ChangeToday / pp.Closing) * 100.00
			}
		} else {
			fmt.Printf("Failed to get day summary %s as %s - %s\n", symbolType, req.As, err)
		}

		resp.Prices[i] = pp
	}

	return resp, nil
}
