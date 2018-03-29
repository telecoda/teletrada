package domain

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

var DefaultArchive = NewSymbolsArchive()

const (
	BTC  = "BTC"
	BNB  = "BNB"
	ETH  = "ETH"
	LTC  = "LTC"
	USDT = "USDT"
)

type SymbolsArchive interface {
	AddSymbol(symbol Symbol) bool
	AddPrice(price Price) error
	GetSymbol(symbol SymbolType) (Symbol, error)
	GetSymbolTypes() map[SymbolType][]SymbolType
	GetLatestPriceAs(base SymbolType, as SymbolType) (Price, error)
	GetPriceAs(base SymbolType, as SymbolType, at time.Time) (Price, error)
	GetDaySummaryAs(base SymbolType, as SymbolType) (DaySummary, error)

	UpdatePrices() error
	UpdateDaySummaries() error

	GetStatus() ArchiveStatus
	// Loading history
	LoadPrices(path string) error
}

type symbolsArchive struct {
	sync.RWMutex
	symbols       map[SymbolType]Symbol
	updateStarted time.Time
	lastUpdated   time.Time
	updateCount   int
	// price persistence
	persist bool
}

type ArchiveStatus struct {
	LastUpdated  time.Time
	UpdateCount  int
	TotalSymbols int
}

func NewSymbolsArchive() SymbolsArchive {
	sa := &symbolsArchive{
		symbols: make(map[SymbolType]Symbol),
	}
	return sa
}

func (sa *symbolsArchive) GetSymbol(symbol SymbolType) (Symbol, error) {
	sa.RLock()
	defer sa.RUnlock()

	if s, ok := sa.symbols[symbol]; ok {
		return s, nil
	}
	return nil, fmt.Errorf("Symbol: %s not found", symbol)
}

func (sa *symbolsArchive) GetSymbolTypes() map[SymbolType][]SymbolType {
	sa.RLock()
	defer sa.RUnlock()

	symbolTypes := make(map[SymbolType][]SymbolType, len(sa.symbols))

	for k, v := range sa.symbols {
		symbolTypes[k] = v.GetAsTypes()
	}

	return symbolTypes
}

// AddSymbol - adds a new symbol to the archive
// returns true is this is actually a New symbol
func (sa *symbolsArchive) AddSymbol(symbol Symbol) bool {
	sa.Lock()
	defer sa.Unlock()
	if _, ok := sa.symbols[symbol.GetType()]; !ok {
		sa.symbols[symbol.GetType()] = symbol
		return true
	}
	return false
}

func (s *symbolsArchive) initPrices() {
}

// GetLatestPriceAs - returns the latest price of base symbol as another symbol
func (sa *symbolsArchive) GetLatestPriceAs(base SymbolType, as SymbolType) (Price, error) {

	price, err := sa.getLatestPriceAs(base, as)
	if err == nil {
		return price, nil
	}

	// no price found for trading pair of base/as
	// so we'll have to convert via BTC

	/* fetching strategy
	base -> BTC (Always fetch BTC price first)
	BTC -> As
	*/
	baseToBtc, err := sa.getLatestPriceAs(base, BTC)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as there is no %s/%s prices", base, as, base, BTC)
	}

	// now get BTC -> as price
	btcToAs, err := sa.getLatestPriceAs(BTC, as)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as there is no %s/%s prices", base, as, BTC, as)
	}

	// combine price conversions for overall exchange rate
	combinedPrice := Price{
		Base:     base,
		As:       as,
		Price:    baseToBtc.Price * btcToAs.Price,
		At:       btcToAs.At,
		Exchange: baseToBtc.Exchange,
	}

	return combinedPrice, nil

}

// getLatestPriceAs - fetches symbol and latest price for it
func (sa *symbolsArchive) getLatestPriceAs(base SymbolType, as SymbolType) (Price, error) {
	// Get symbol
	baseSymbol, err := sa.GetSymbol(base)
	if err != nil {
		return Price{}, fmt.Errorf("No prices for symbol %q", base)
	}

	price, err := baseSymbol.GetLatestPriceAs(as)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as their is no %s/%s prices", base, as, base, as)
	}
	return price, nil
}

// GetPriceAs - returns the price of base symbol as another symbol at a particular time
func (sa *symbolsArchive) GetPriceAs(base SymbolType, as SymbolType, at time.Time) (Price, error) {

	price, err := sa.getPriceAs(base, as, at)
	if err == nil {
		return price, nil
	}

	// no price found for trading pair of base/as
	// so we'll have to convert via BTC

	/* fetching strategy
	base -> BTC (Always fetch BTC price first)
	BTC -> As
	*/
	baseToBtc, err := sa.getPriceAs(base, BTC, at)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as there is no %s/%s prices at %s", base, as, base, BTC, at.Format(DATE_FORMAT))
	}

	// now get BTC -> as price
	btcToAs, err := sa.getPriceAs(BTC, as, at)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as there is no %s/%s prices at %s", base, as, BTC, as, at.Format(DATE_FORMAT))
	}

	// combine price conversions for overall exchange rate
	combinedPrice := Price{
		Base:     base,
		As:       as,
		Price:    baseToBtc.Price * btcToAs.Price,
		At:       at,
		Exchange: baseToBtc.Exchange,
	}

	return combinedPrice, nil

}

// getPriceAs - fetches symbol and latest price for it at a particular time
func (sa *symbolsArchive) getPriceAs(base SymbolType, as SymbolType, at time.Time) (Price, error) {
	// Get symbol
	baseSymbol, err := sa.GetSymbol(base)
	if err != nil {
		return Price{}, fmt.Errorf("No prices for symbol %q", base)
	}

	price, err := baseSymbol.GetPriceAs(as, at)
	if err != nil {
		return Price{}, fmt.Errorf("unable to convert %q to %q as their is no %s/%s prices at %s", base, as, base, as, at.Format(DATE_FORMAT))
	}
	return price, nil
}

// GetDaySummaryAs - returns the last days summary of base symbol as another symbol
func (sa *symbolsArchive) GetDaySummaryAs(base SymbolType, as SymbolType) (DaySummary, error) {

	baseSymbol, err := sa.GetSymbol(base)
	if err != nil {
		return DaySummary{}, fmt.Errorf("No prices for symbol %q", base)
	}

	sum, err := baseSymbol.GetDaySummaryAs(as)
	if err != nil {
		return DaySummary{}, fmt.Errorf("no day summary for %q as %q", base, as)
	}
	return sum, nil

}

func (sa *symbolsArchive) UpdatePrices() error {

	exPrices, err := DefaultClient.GetLatestPrices()
	if err != nil {
		return fmt.Errorf("Failed to get latest prices: %s", err)
	}

	prices := make([]Price, len(exPrices))

	for i, exPrice := range exPrices {
		// convert Exchange price to Domain price
		prices[i] = Price{
			Base:     SymbolType(exPrice.Base),
			As:       SymbolType(exPrice.As),
			Price:    exPrice.Price,
			At:       exPrice.At,
			Exchange: exPrice.Exchange,
		}
	}
	// process latest prices
	for _, price := range prices {
		if err := sa.savePrice(price); err != nil {
			return err
		}
	}

	// send to influxDB
	if err := DefaultMetrics.SavePriceMetrics(prices); err != nil {
		return err
	}

	sa.Lock()
	sa.updateCount++
	sa.lastUpdated = time.Now()
	sa.Unlock()
	return nil
}

func (sa *symbolsArchive) AddPrice(price Price) error {
	return sa.savePrice(price)
}

func (sa *symbolsArchive) UpdateDaySummaries() error {
	summaries, err := DefaultClient.GetDaySummaries()
	if err != nil {
		return err
	}
	fmt.Printf("Starting update daily summaries %d\n", len(summaries))
	for _, exSummary := range summaries {
		symbol, err := sa.GetSymbol(SymbolType(exSummary.Base))
		if err != nil {
			log.Printf("ERROR: [UpdateDaySummaries] failed getting symbol %s - %s", exSummary.Base, err)
			continue
		}

		summary := DaySummary{
			Base:             SymbolType(exSummary.Base),
			As:               SymbolType(exSummary.As),
			OpenPrice:        exSummary.OpenPrice,
			ClosePrice:       exSummary.ClosePrice,
			WeightedAvgPrice: exSummary.WeightedAvgPrice,
			HighestPrice:     exSummary.HighestPrice,
			LowestPrice:      exSummary.LowestPrice,
			ChangePrice:      exSummary.ChangePrice,
			ChangePercent:    exSummary.ChangePercent,
			At:               exSummary.At,
			Exchange:         exSummary.Exchange,
		}

		symbol.AddDaySummary(summary)
	}

	return nil
}

func (sa *symbolsArchive) GetStatus() ArchiveStatus {
	sa.RLock()
	defer sa.RUnlock()

	return ArchiveStatus{
		LastUpdated:  sa.lastUpdated,
		UpdateCount:  sa.updateCount,
		TotalSymbols: len(sa.symbols),
	}
}

// savePrice - saves a price in the archive and updates latest if this is most recent
func (sa *symbolsArchive) savePrice(price Price) error {

	if err := price.Validate(); err != nil {
		return fmt.Errorf("Price is not valid: %s - %#v", err, price)
	}

	var pSymbol Symbol
	var err error

	pSymbol, err = sa.GetSymbol(price.Base)

	if err != nil {
		// create new symbol
		pSymbol = NewSymbol(price.Base)
		// add to map
		sa.AddSymbol(pSymbol)
		log.Printf("New Symbol added: %s\n", pSymbol.GetType())
	}
	pSymbol.AddPrice(price)

	return nil
}

/*

Metrics we want to save:

Every minute:
	- All coin prices in following currencies
		BTC
		ETH
		USDT
		GBP

Format:

	Point: coin_price
	Tags:
		symbol : coin
	Fields:
		price.BTC - price
		price.ETH - price
		price.USDT - price
		price.GBP - price

*/

func (sa *symbolsArchive) LoadPrices(dir string) error {
	// check dir exists
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Can't load from dir: %s - %s", dir, err)
	}

	fmt.Printf("Loading prices from %d files\n", len(files))
	t := time.Now()
	// read all files
	for _, file := range files {
		// only load from .json files
		if strings.HasSuffix(file.Name(), ".json") {
			if err := sa.loadPricesFrom(filepath.Join(dir, file.Name())); err != nil {
				return err
			}
		}
	}
	fmt.Printf("Loaded in %s\n", time.Now().Sub(t).String())

	return nil
}

func (sa *symbolsArchive) loadPricesFrom(filePath string) error {

	f, err := os.OpenFile(filePath, 0, 0)
	if err != nil {
		return err
	}
	defer f.Close()
	b := bufio.NewReader(f)

	pricesJSON, err := ioutil.ReadAll(b)
	if err != nil {
		return err
	}

	prices := make([]Price, 0)
	err = json.Unmarshal(pricesJSON, &prices)
	if err != nil {
		return err
	}

	for _, price := range prices {
		if err := sa.savePrice(price); err != nil {
			return fmt.Errorf("Failed to load price from file: %s - %s", filePath, err)
		}
	}

	return nil
}
