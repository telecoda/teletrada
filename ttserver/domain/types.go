package domain

import (
	"time"
)

const (
	USDT = "USDT"
)

type archive struct {
}

type tradeHistory struct {
	trades []trade
}

type trade struct {
	symbol    string
	quantity  float64
	exchange  string
	price     float64
	fee       float64
	totalCost float64
	date      time.Time
}

type buy trade

type sell trade
