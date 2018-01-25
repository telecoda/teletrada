package domain

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/exchanges"
)

func TestAddSymbol(t *testing.T) {
	testSymbol := SymbolType("tester")
	unknownSymbol := SymbolType("unknown")

	symbol := NewSymbol(testSymbol)

	archive := NewSymbolsArchive()

	_, err := archive.GetSymbol(unknownSymbol)
	assert.Error(t, err, "Should return an erro")

	newSymbol := archive.AddSymbol(symbol)
	returnedSymbol, err := archive.GetSymbol(testSymbol)
	assert.NoError(t, err, "Should not return an error")
	assert.NotNil(t, returnedSymbol, "Should not be nil")
	assert.True(t, newSymbol, "A new symbol has been added")

	notAddedSymbol := archive.AddSymbol(symbol)
	assert.False(t, notAddedSymbol, "A new symbol has not been added")

}

func TestSavePrice(t *testing.T) {
	testSymbol := SymbolType("tester")
	priceSymbol := SymbolType(USDT)
	priceYest := 1234.56
	priceToday := 2345.67
	today := time.Now()
	yesterday := today.AddDate(0, 0, -1)

	archive := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}

	yesterdaysPrice := Price{
		Base:  testSymbol,
		As:    priceSymbol,
		Price: priceYest,
		At:    yesterday,
	}

	// saving a price should automatically add a symbol
	archive.savePrice(yesterdaysPrice)

	// get price afterwards
	symbol, err := archive.GetSymbol(testSymbol)
	assert.NoError(t, err, "Should not return an error")
	assert.NotNil(t, symbol, "Should not be nil")

	latestPrice, err := symbol.GetLatestPriceAs(USDT)
	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, priceYest, latestPrice.Price, "Should match yesterday's price")

	// Add another price
	todaysPrice := Price{
		Base:  testSymbol,
		As:    priceSymbol,
		Price: priceToday,
		At:    today,
	}

	// fetch new referesh to same price
	newSymbol, err := archive.GetSymbol(testSymbol)
	assert.NoError(t, err, "Should not return an error")

	archive.savePrice(todaysPrice)

	// Price should be updated on both references to same symbol
	newLatestPrice, err := newSymbol.GetLatestPriceAs(USDT)
	assert.Equal(t, priceToday, newLatestPrice.Price, "Should match today's price")

	// without refetching symbol we should have the latest price
	latestPrice, err = symbol.GetLatestPriceAs(USDT)

	assert.NoError(t, err, "Should not return an error")
	assert.Equal(t, priceToday, latestPrice.Price, "Should match today's price")

}

func TestUpdatePrices(t *testing.T) {

	archive := NewSymbolsArchive()

	mc, err := exchanges.NewMockClient()
	assert.NoError(t, err)

	// run test with mocked data
	DefaultClient = mc

	testSymbol := SymbolType("BTC")
	currency := SymbolType("USDT")

	_, err = archive.GetSymbol(testSymbol)
	assert.Error(t, err, "Error expected no symbols loaded")

	err = archive.UpdatePrices()
	assert.NoError(t, err, "Should not return an error")

	btc, err := archive.GetSymbol(testSymbol)
	assert.NoError(t, err)
	assert.Equal(t, btc.GetType(), testSymbol)

	latestPrice, err := btc.GetLatestPriceAs(currency)
	assert.NoError(t, err)
	assert.Equal(t, testSymbol, latestPrice.Base)
	assert.Equal(t, currency, latestPrice.As)
	assert.Equal(t, 12000.12345, latestPrice.Price)

}

func TestScheduledUpdate(t *testing.T) {

	archive := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}

	mc, err := exchanges.NewMockClient()
	assert.NoError(t, err)

	// run test with mocked data
	DefaultClient = mc

	assert.Equal(t, 0, archive.updateCount, "No updates yet")

	freq := time.Duration(100 * time.Millisecond)

	archive.StartUpdater(freq)
	time.Sleep(550 * time.Millisecond)

	assert.Equal(t, 5, archive.updateCount, "5 updates should have happened by now")

}

func TestPricePersistence(t *testing.T) {

	archive := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}

	mc, err := exchanges.NewMockClient()
	assert.NoError(t, err)

	// run test with mocked data
	DefaultClient = mc

	testDir := filepath.Join(".", "testPriceHistory")

	// clear test files
	info, err := ioutil.ReadDir(testDir)
	assert.NoError(t, err)

	for _, file := range info {
		err = os.Remove(filepath.Join(testDir, file.Name()))
		assert.NoError(t, err)
	}

	archive.StartPersistence(testDir)

	err = archive.UpdatePrices()
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = archive.UpdatePrices()
	assert.NoError(t, err)

	// check files
	info, err = ioutil.ReadDir(testDir)
	assert.NoError(t, err)

	assert.Equal(t, 2, len(info), "There should be 2 files in the directory now")
}
