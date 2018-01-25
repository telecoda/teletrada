package domain

import (
	"fmt"
	"time"
)

type Price struct {
	Base  SymbolType
	As    SymbolType
	Price float64
	At    time.Time
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
