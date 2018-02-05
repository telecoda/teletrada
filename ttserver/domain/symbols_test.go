package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSymbolPrices(t *testing.T) {
	testSymbol := SymbolType("tester")
	symbol := NewSymbol(testSymbol)

	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	// add today's USD price
	priceToday := Price{
		Base:  testSymbol,
		As:    USDT,
		Price: 20000.00,
		At:    today,
	}

	// add yesterday;s USD price
	// this checks that prices are sorted when they are added
	priceYesterday := Price{
		Base:  testSymbol,
		As:    USDT,
		Price: 10000.00,
		At:    yesterday,
	}

	symbol.AddPrice(priceToday)
	symbol.AddPrice(priceYesterday)

	// get latest
	latestUSDPrice, err := symbol.GetLatestPriceAs(USDT)

	assert.NoError(t, err)
	assert.Equal(t, 20000.00, latestUSDPrice.Price, "Failed to get latest price")

	// get todays price
	todaysUSDPrice, err := symbol.GetPriceAs(USDT, today)

	assert.NoError(t, err)
	assert.Equal(t, 20000.00, todaysUSDPrice.Price, "Failed to get today's price")

	// get yesterdays price
	yesterdaysUSDPrice, err := symbol.GetPriceAs(USDT, yesterday)

	assert.NoError(t, err)
	assert.Equal(t, 10000.00, yesterdaysUSDPrice.Price, "Failed to get yesterday's price")

	// get price inbetween two dates
	// it should calculate the mid price between the tow
	halfDay := today.Add(-12 * time.Hour)
	halfdayUSDPrice, err := symbol.GetPriceAs(USDT, halfDay)
	assert.NoError(t, err)
	assert.Equal(t, 15000.00, halfdayUSDPrice.Price, "failed to get an adjusted price")

	// get price before historic prices exist
	beforeDate := yesterday.Add(-12 * time.Hour)
	beforeUSDPrice, err := symbol.GetPriceAs(USDT, beforeDate)
	assert.NoError(t, err)
	assert.Equal(t, 10000.00, beforeUSDPrice.Price, "failed to get a price")

	// get unknown symbol
	unknown := SymbolType("unknown")
	_, err = symbol.GetLatestPriceAs(unknown)
	assert.Error(t, err, "Should return an error")

	_, err = symbol.GetPriceAs(unknown, today)
	assert.Error(t, err, "Should return an error")

}

