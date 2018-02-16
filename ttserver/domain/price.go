package domain

import (
	"fmt"
	"time"
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
