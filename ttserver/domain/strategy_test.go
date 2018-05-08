package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

func strategySetup(t *testing.T) {
	err := initMockClients(nil)
	assert.NoError(t, err)
}
func TestBaseStrategy(t *testing.T) {
	symbol := SymbolType("test-symbol")
	as := SymbolType("USDT")

	strat, err := NewBaseStrategy("base", symbol, as, 100.00)
	assert.NoError(t, err)
	assert.False(t, strat.IsRunning())
	strat.Start()
	assert.True(t, strat.IsRunning())
	strat.Stop()
	assert.False(t, strat.IsRunning())

	now := servertime.Now()

	met, err := strat.ConditionMet(now)
	assert.False(t, met)
	assert.NoError(t, err)

}
func TestDoNothingStrategy(t *testing.T) {
	symbol := SymbolType("test-symbol")
	as := SymbolType("USDT")

	strat, err := NewDoNothingStrategy("do-nothing", symbol, as, 100.00)
	assert.NoError(t, err)

	assert.False(t, strat.IsRunning())
	strat.Start()
	assert.True(t, strat.IsRunning())
	strat.Stop()
	assert.False(t, strat.IsRunning())

	now := servertime.Now()

	met, err := strat.ConditionMet(now)
	assert.False(t, met)
	assert.NoError(t, err)

}

func TestPriceAboveStrategy(t *testing.T) {

	servertime.UseFakeTime()
	defer servertime.UseRealTime()

	strategySetup(t)

	today := servertime.Now()
	tomorrow := today.AddDate(0, 0, 1)
	yesterday := today.AddDate(0, 0, -1)

	symbol := SymbolType("test-symbol")
	as := SymbolType("USDT")

	abovePrice := 200.00
	belowPrice := 50.00
	coinPercent := 100.00

	tests := []struct {
		name         string
		price        Price
		conditionMet bool
	}{
		{
			name:         "Before price change",
			price:        Price{Base: symbol, As: as, Price: 100.00, At: yesterday, Exchange: "exchange"},
			conditionMet: false,
		},
		{
			name:         "Price above - sell!",
			price:        Price{Base: symbol, As: as, Price: abovePrice + 1, At: today, Exchange: "exchange"},
			conditionMet: true,
		},
		{
			name:         "Price below - do nothing",
			price:        Price{Base: symbol, As: as, Price: belowPrice - 1, At: tomorrow, Exchange: "exchange"},
			conditionMet: false,
		},
	}

	strat, err := NewPriceAboveStrategy("sell-price", symbol, as, abovePrice, coinPercent)
	assert.NoError(t, err)

	assert.False(t, strat.IsRunning())
	strat.Start()
	assert.True(t, strat.IsRunning())
	strat.Stop()
	assert.False(t, strat.IsRunning())

	strat.Start()

	// save prices before processing
	for _, test := range tests {
		err := DefaultArchive.AddPrice(test.price)
		assert.NoError(t, err)
	}

	for _, test := range tests {
		met, err := strat.ConditionMet(test.price.At)
		assert.NoError(t, err)
		assert.Equal(t, test.conditionMet, met)
	}

}

func TestPriceBelowStrategy(t *testing.T) {

	strategySetup(t)

	today := servertime.Now()
	tomorrow := today.AddDate(0, 0, 1)
	yesterday := today.AddDate(0, 0, -1)

	symbol := SymbolType("test-symbol")
	as := SymbolType("USDT")

	abovePrice := 200.00
	belowPrice := 50.00
	coinPercent := 100.00

	tests := []struct {
		name         string
		price        Price
		conditionMet bool
	}{
		{
			name:         "Before price change",
			price:        Price{Base: symbol, As: as, Price: 100.00, At: yesterday, Exchange: "exchange"},
			conditionMet: false,
		},
		{
			name:         "Price above - sell!",
			price:        Price{Base: symbol, As: as, Price: abovePrice + 1, At: today, Exchange: "exchange"},
			conditionMet: false,
		},
		{
			name:         "Price below - do nothing",
			price:        Price{Base: symbol, As: as, Price: belowPrice - 1, At: tomorrow, Exchange: "exchange"},
			conditionMet: true,
		},
	}

	strat, err := NewPriceBelowStrategy("buy-price", symbol, as, belowPrice, coinPercent)
	assert.NoError(t, err)

	assert.False(t, strat.IsRunning())
	strat.Start()
	assert.True(t, strat.IsRunning())
	strat.Stop()
	assert.False(t, strat.IsRunning())

	strat.Start()

	// save prices before processing
	for _, test := range tests {
		err := DefaultArchive.AddPrice(test.price)
		assert.NoError(t, err)
	}

	for _, test := range tests {
		met, err := strat.ConditionMet(test.price.At)
		assert.NoError(t, err)
		assert.Equal(t, test.conditionMet, met)
	}

}
