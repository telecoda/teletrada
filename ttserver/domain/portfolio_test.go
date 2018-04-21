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

	clonedPortfolio, err := originalPortfolio.clone()
	assert.NoError(t, err)

	assert.Equal(t, len(originalPortfolio.balances), len(clonedPortfolio.balances))
	assert.Equal(t, "Live[cloned]", clonedPortfolio.name)
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

func TestDiffNoDifferences(t *testing.T) {

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

	clonedPortfolio, err := originalPortfolio.clone()
	assert.NoError(t, err)

	// Now calc diff (should be zero)
	diff, err := clonedPortfolio.diff(originalPortfolio)
	assert.NoError(t, err)

	assert.Equal(t, len(originalPortfolio.balances), len(diff.balances))
	assert.Equal(t, "Live[cloned][diff]", diff.name)
	assert.False(t, diff.isLive)

	for symbol, balance := range originalPortfolio.balances {
		assert.Equal(t, balance.Symbol, diff.balances[symbol].Symbol)
		assert.Equal(t, balance.Exchange, diff.balances[symbol].Exchange)
		assert.Zero(t, diff.balances[symbol].Free)
		assert.Zero(t, diff.balances[symbol].Locked)
		assert.Zero(t, diff.balances[symbol].Total)
		assert.Equal(t, balance.At, diff.balances[symbol].At)
		assert.Equal(t, balance.As, diff.balances[symbol].As)
		assert.Zero(t, diff.balances[symbol].Price)
		assert.Zero(t, diff.balances[symbol].Value)
		assert.Zero(t, diff.balances[symbol].Price24H)
		assert.Zero(t, diff.balances[symbol].Value24H)
		assert.Zero(t, diff.balances[symbol].Change24H)
		assert.Zero(t, diff.balances[symbol].ChangePct24H)
	}

	// update original
	originalPortfolio.balances["symbol1"].Total = 99

	// make sure cloned value not changed
	assert.Equal(t, float64(99), originalPortfolio.balances["symbol1"].Total)
	assert.Equal(t, float64(0), diff.balances["symbol1"].Total)

}

func TestDiffWithDifferences(t *testing.T) {

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

	clonedPortfolio, err := originalPortfolio.clone()
	assert.NoError(t, err)

	// Now amend cloned portfolio a bit

	symbol1, ok := clonedPortfolio.balances["symbol1"]
	assert.True(t, ok)
	symbol1.Free += 5
	symbol1.Locked += 10
	symbol1.Total += 15
	symbol1.Price += 20
	symbol1.Price24H += 25
	symbol1.Value += 30
	symbol1.Value24H += 35
	symbol1.Change24H += 40
	symbol1.ChangePct24H += 0.1

	clonedPortfolio.balances["symbol1"] = symbol1

	// Now calc diff (should not be zero)
	diff, err := clonedPortfolio.diff(originalPortfolio)
	assert.NoError(t, err)
	assert.Equal(t, len(originalPortfolio.balances), len(diff.balances))
	assert.Equal(t, "Live[cloned][diff]", diff.name)
	assert.False(t, diff.isLive)

	// check symbol1 has changed

	ob, ok := originalPortfolio.balances["symbol1"]
	// check were down with obb (90s hip hop reference for the oldies)
	assert.True(t, ok)
	db, ok := diff.balances["symbol1"]
	assert.True(t, ok)
	assert.Equal(t, ob.Symbol, db.Symbol)
	assert.Equal(t, ob.Exchange, db.Exchange)
	assert.Equal(t, float64(5), db.Free)
	assert.Equal(t, float64(10), db.Locked)
	assert.Equal(t, ob.At, db.At)
	assert.Equal(t, ob.As, db.As)
	assert.Equal(t, float64(15), db.Total)
	assert.Equal(t, float64(20), db.Price)
	assert.Equal(t, float64(25), db.Price24H)
	assert.Equal(t, float64(30), db.Value)
	assert.Equal(t, float64(35), db.Value24H)
	assert.Equal(t, float64(40), db.Change24H)
	assert.InDelta(t, float64(0.1), db.ChangePct24H, 0.000000001)

}
