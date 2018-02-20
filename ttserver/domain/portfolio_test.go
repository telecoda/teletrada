package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/exchanges"
)

func TestClone(t *testing.T) {

	originalPortfolio := &portfolio{
		name:     "Live",
		isLive:   true,
		balances: make(map[SymbolType]*BalanceAs),
	}

	now := time.Now()
	// add some balances
	balance1 := &BalanceAs{
		CoinBalance: exchanges.CoinBalance{
			Symbol:   "symbol1",
			Exchange: "test-exchange",
			Free:     5,
			Locked:   5,
		},
		Total:        10,
		At:           now,
		As:           SymbolType("BTC"),
		Price:        25.00,
		Value:        250.00,
		Price24H:     20.00,
		Value24H:     200.00,
		Change24H:    50.00,
		ChangePct24H: 25.0,
	}
	balance2 := &BalanceAs{
		CoinBalance: exchanges.CoinBalance{
			Symbol:   "symbol2",
			Exchange: "test-exchange",
			Free:     50,
			Locked:   50,
		},
		Total:        10,
		At:           now,
		As:           SymbolType("BTC"),
		Price:        25.00,
		Value:        2500.00,
		Price24H:     20.00,
		Value24H:     2000.00,
		Change24H:    500.00,
		ChangePct24H: 25.0,
	}

	originalPortfolio.balances["symbol1"] = balance1
	originalPortfolio.balances["symbol2"] = balance2

	clonedPortfolio, err := originalPortfolio.clone("Cloned portfolio")
	assert.NoError(t, err)

	assert.Equal(t, len(originalPortfolio.balances), len(clonedPortfolio.balances))
	assert.Equal(t, "Cloned portfolio", clonedPortfolio.name)
	assert.False(t, clonedPortfolio.isLive)

	for symbol, balance := range originalPortfolio.balances {
		assert.Equal(t, balance, clonedPortfolio.balances[symbol])
	}

	// update original
	originalPortfolio.balances["symbol1"].Total = 99

	// make sure cloned value not changed
	assert.Equal(t, float64(99), originalPortfolio.balances["symbol1"].Total)
	assert.Equal(t, float64(10), clonedPortfolio.balances["symbol1"].Total)

}
