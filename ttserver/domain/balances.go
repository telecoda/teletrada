package domain

import (
	"sync"
	"time"

	"github.com/telecoda/teletrada/exchanges"
)

type BalanceAs struct {
	sync.RWMutex
	exchanges.CoinBalance
	Total        float64
	At           time.Time
	As           SymbolType
	Price        float64
	Value        float64
	Price24H     float64
	Value24H     float64
	Change24H    float64
	ChangePct24H float64
	//
	BuyStrategy  Strategy
	SellStrategy Strategy
}
