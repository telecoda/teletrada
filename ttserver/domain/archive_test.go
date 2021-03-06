package domain

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

func setupArchive() *symbolsArchive {

	// TODO drop / create test influx db here...

	archive := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}

	return archive
}

func NewTestSymbolsArchive() SymbolsArchive {
	sa := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}
	return sa
}

func TestAddSymbol(t *testing.T) {
	testSymbol := SymbolType("tester")
	unknownSymbol := SymbolType("unknown")

	symbol := NewSymbol(testSymbol)

	archive := NewTestSymbolsArchive()

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
	today := servertime.Now()
	yesterday := today.AddDate(0, 0, -1)

	archive := setupArchive()

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
		Base:     testSymbol,
		As:       priceSymbol,
		Price:    priceToday,
		At:       today,
		Exchange: "test_exchange",
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

	_, err := initMockServer()
	assert.NoError(t, err)

	archive := NewTestSymbolsArchive()

	testSymbol := SymbolType(BTC)
	currency := SymbolType(USDT)

	_, err = archive.GetSymbol(testSymbol)
	assert.Error(t, err, "Error expected no symbols loaded")

	err = archive.UpdatePrices()
	assert.NoError(t, err, "Should not return an error")

	btc, err := archive.GetSymbol(testSymbol)
	if assert.NoError(t, err) {
		assert.Equal(t, btc.GetType(), testSymbol)

		latestPrice, err := btc.GetLatestPriceAs(currency)
		assert.NoError(t, err)
		assert.Equal(t, testSymbol, latestPrice.Base)
		assert.Equal(t, currency, latestPrice.As)
		assert.Equal(t, _btcAsUsdt, latestPrice.Price)
	}

}

func TestPricePersistence(t *testing.T) {

	archive := setupArchive()

	_, err := initMockServer()
	assert.NoError(t, err)

	err = archive.UpdatePrices()
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)
	err = archive.UpdatePrices()
	assert.NoError(t, err)

	// read data from influx
	// TODO
}

func TestMultiCurrencyPrices(t *testing.T) {

	/* this test is to check we can convert prices to different currencies
	 */

	archive := setupArchive()

	_, err := initMockServer()
	assert.NoError(t, err)

	ltcSymbol := SymbolType("LTC")
	btcSymbol := SymbolType("BTC")
	ethSymbol := SymbolType("ETH")
	usdtSymbol := SymbolType("USDT")

	today := servertime.Now()

	// add LTC -> BTC price
	LtcBtcPrice := Price{
		Base:     ltcSymbol,
		As:       btcSymbol,
		Price:    0.1, // how much 1 LTC is worth in BTC
		At:       today,
		Exchange: "test_exchange",
	}

	// add BTC -> ETH price
	BtcEthPrice := Price{
		Base:     btcSymbol,
		As:       ethSymbol,
		Price:    20.0, // how much 1 BTC is worth in ETH
		At:       today,
		Exchange: "test_exchange",
	}

	// add BTC -> USDT price
	BtcUsdtPrice := Price{
		Base:     btcSymbol,
		As:       usdtSymbol,
		Price:    20000.0, // how much 1 BTC is worth in USDT
		At:       today,
		Exchange: "test_exchange",
	}

	err = archive.savePrice(LtcBtcPrice)
	assert.NoError(t, err)

	err = archive.savePrice(BtcEthPrice)
	assert.NoError(t, err)

	err = archive.savePrice(BtcUsdtPrice)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		base     SymbolType
		as       SymbolType
		expErr   error
		expPrice float64
	}{
		{
			name:     "Simple conversion",
			base:     ltcSymbol,
			as:       btcSymbol,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price,
		},
		{
			name:     "LTC -> ETH (via BTC)",
			base:     ltcSymbol,
			as:       ethSymbol,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price * BtcEthPrice.Price,
		},
		{
			name:     "LTC -> USDT (via BTC)",
			base:     ltcSymbol,
			as:       usdtSymbol,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price * BtcUsdtPrice.Price,
		},
		{
			name:   "Unknown base symbol",
			base:   SymbolType("unknown"),
			as:     usdtSymbol,
			expErr: fmt.Errorf(`unable to convert "unknown" to "USDT" as there is no unknown/BTC prices`),
		},
		{
			name:   "Unknown as symbol",
			base:   ltcSymbol,
			as:     SymbolType("unknown"),
			expErr: fmt.Errorf(`unable to convert "LTC" to "unknown" as there is no BTC/unknown prices`),
		},
		{
			name:   "No BTC prices",
			base:   ethSymbol,
			as:     ltcSymbol,
			expErr: fmt.Errorf(`unable to convert "ETH" to "LTC" as there is no ETH/BTC prices`),
		},
	}

	for _, test := range tests {
		price, err := archive.GetLatestPriceAs(test.base, test.as)

		if test.expErr == nil && err != nil {
			assert.Fail(t, fmt.Sprintf("Didn't expect error and received %s", err), "Test %s", test.name)
		}

		if test.expErr != nil && err == nil {
			assert.Fail(t, fmt.Sprintf("Expected error %s and didn't receive it", test.expErr), "Test %s", test.name)
		}

		if test.expErr != nil && err != nil {
			assert.Equal(t, test.expErr, err)
		}

		// if no error compare result
		if test.expErr == nil {
			assert.Equal(t, test.expPrice, price.Price, "Test %s", test.name)
		}

	}

}

func TestMultiCurrencyPricesAt(t *testing.T) {

	/* this test is to check we can convert prices to different currencies
	 */

	archive := setupArchive()

	_, err := initMockServer()
	assert.NoError(t, err)

	ltcSymbol := SymbolType("LTC")
	btcSymbol := SymbolType("BTC")
	ethSymbol := SymbolType("ETH")
	usdtSymbol := SymbolType("USDT")

	today := servertime.Now()
	yesterday := today.AddDate(0, 0, -1)

	// add LTC -> BTC price
	LtcBtcPrice := Price{
		Base:     ltcSymbol,
		As:       btcSymbol,
		Price:    0.1, // how much 1 LTC is worth in BTC
		At:       today,
		Exchange: "test_exchange",
	}

	// add BTC -> ETH price
	BtcEthPrice := Price{
		Base:     btcSymbol,
		As:       ethSymbol,
		Price:    20.0, // how much 1 BTC is worth in ETH
		At:       today,
		Exchange: "test_exchange",
	}

	// add BTC -> USDT price
	BtcUsdtPrice := Price{
		Base:     btcSymbol,
		As:       usdtSymbol,
		Price:    20000.0, // how much 1 BTC is worth in USDT
		At:       today,
		Exchange: "test_exchange",
	}

	// add BTC -> USDT price (yesterdays price)
	BtcUsdtPriceYesterday := Price{
		Base:     btcSymbol,
		As:       usdtSymbol,
		Price:    10000.0, // how much 1 BTC is worth in USDT
		At:       yesterday,
		Exchange: "test_exchange",
	}

	err = archive.savePrice(LtcBtcPrice)
	assert.NoError(t, err)

	err = archive.savePrice(BtcEthPrice)
	assert.NoError(t, err)

	err = archive.savePrice(BtcUsdtPrice)
	assert.NoError(t, err)

	err = archive.savePrice(BtcUsdtPriceYesterday)
	assert.NoError(t, err)

	tests := []struct {
		name     string
		base     SymbolType
		as       SymbolType
		at       time.Time
		expErr   error
		expPrice float64
	}{
		{
			name:     "Simple conversion",
			base:     ltcSymbol,
			as:       btcSymbol,
			at:       today,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price,
		},
		{
			name:     "LTC -> ETH (via BTC)",
			base:     ltcSymbol,
			as:       ethSymbol,
			at:       today,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price * BtcEthPrice.Price,
		},
		{
			name:     "LTC -> USDT (via BTC)",
			base:     ltcSymbol,
			as:       usdtSymbol,
			at:       today,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price * BtcUsdtPrice.Price,
		},
		{
			name:     "LTC -> USDT (via BTC) for yesterday",
			base:     ltcSymbol,
			as:       usdtSymbol,
			at:       yesterday,
			expErr:   nil,
			expPrice: LtcBtcPrice.Price * BtcUsdtPriceYesterday.Price,
		},
		{
			name:   "Unknown base symbol",
			base:   SymbolType("unknown"),
			as:     usdtSymbol,
			at:     today,
			expErr: fmt.Errorf(`unable to convert "unknown" to "USDT" as there is no unknown/BTC prices at %s`, today.Format(DATE_FORMAT)),
		},
		{
			name:   "Unknown as symbol",
			base:   ltcSymbol,
			as:     SymbolType("unknown"),
			at:     today,
			expErr: fmt.Errorf(`unable to convert "LTC" to "unknown" as there is no BTC/unknown prices at %s`, today.Format(DATE_FORMAT)),
		},
		{
			name:   "No BTC prices",
			base:   ethSymbol,
			as:     ltcSymbol,
			at:     today,
			expErr: fmt.Errorf(`unable to convert "ETH" to "LTC" as there is no ETH/BTC prices at %s`, today.Format(DATE_FORMAT)),
		},
	}

	for _, test := range tests {
		price, err := archive.GetPriceAs(test.base, test.as, test.at)

		if test.expErr == nil && err != nil {
			assert.Fail(t, fmt.Sprintf("Didn't expect error and received %s", err), "Test %s", test.name)
		}

		if test.expErr != nil && err == nil {
			assert.Fail(t, fmt.Sprintf("Expected error %s and didn't receive it", test.expErr), "Test %s", test.name)
		}

		if test.expErr != nil && err != nil {
			assert.Equal(t, test.expErr, err)
		}

		// if no error compare result
		if test.expErr == nil {
			assert.Equal(t, test.expPrice, price.Price, "Test %s", test.name)
		}

	}

}
