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
	"text/tabwriter"
	"time"
)

var DefaultArchive = NewSymbolsArchive()

type SymbolsArchive interface {
	AddSymbol(symbol Symbol) bool
	GetSymbol(symbol SymbolType) (Symbol, error)
	ListPrices(incHistory bool) // include historic prices
	UpdatePrices() error
	// Starts automatic price updater
	StartUpdater(frequency time.Duration)
	// Stops automatic price updater
	StopUpdater()
	StartPersistence(path string) error
	StopPersistence()
	persistPrices(prices []Price) error
	// Loading history
	LoadPrices(path string) error
}

type symbolsArchive struct {
	sync.RWMutex
	symbols map[SymbolType]Symbol
	// scheduling
	stopUpdate    chan bool
	updateStarted time.Time
	updateCount   int
	// price persistence
	persistToDisk bool
	persistDir    string
}

func NewSymbolsArchive() SymbolsArchive {
	sa := &symbolsArchive{
		symbols:    make(map[SymbolType]Symbol),
		stopUpdate: make(chan bool),
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

func (sa *symbolsArchive) UpdatePrices() error {
	exPrices, err := DefaultClient.GetLatestPrices()
	if err != nil {
		return fmt.Errorf("Failed to get latest prices: %s", err)
	}

	prices := make([]Price, len(exPrices))

	for i, exPrice := range exPrices {
		// convert Exchange price to Domain price
		prices[i] = Price{
			Base:  SymbolType(exPrice.Base),
			As:    SymbolType(exPrice.As),
			Price: exPrice.Price,
			At:    exPrice.At,
		}
	}
	// process latest prices
	for _, price := range prices {
		if err := sa.savePrice(price); err != nil {
			return err
		}
	}

	if sa.persistToDisk {
		if err := sa.persistPrices(prices); err != nil {
			return err
		}
	}

	sa.Lock()
	sa.updateCount++
	sa.Unlock()
	return nil
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

func (sa *symbolsArchive) ListPrices(incHistory bool) {
	sa.RLock()
	defer sa.RUnlock()

	fmt.Printf("Prices\n")
	fmt.Printf("======\n")

	tw := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight)

	// headings
	fmt.Fprintf(tw, "base\tas\tprice\t\t         at\n")
	fmt.Fprintf(tw, "======\t======\t=============\t\t===================\n")

	for _, symbol := range sa.symbols {

		// get USDT prices
		price, err := symbol.GetLatestPriceAs(USDT)
		if err != nil {
			fmt.Printf("Error getting USDT price: %s\n", err)
			continue
		}
		fmt.Fprintf(tw, "(Latest) %s\t%s\t%s\t\t%s\n", price.Base, price.As, fmt.Sprintf("%f", price.Price), price.At.Format(DATE_FORMAT))

		if incHistory {
			// yesterday's prices
			now := time.Now().UTC()
			yesterday := now.AddDate(0, 0, -1)
			price, err := symbol.GetPriceAs(USDT, yesterday)
			if err != nil {
				fmt.Printf("Error getting USDT price: %s\n", err)
				continue
			}
			fmt.Fprintf(tw, "(Yesterday) %s\t%s\t%s\t\t%s\n", price.Base, price.As, fmt.Sprintf("%f", price.Price), price.At.Format(DATE_FORMAT))

		}
	}

	tw.Flush()
}

// Stops automatic price updater
func (sa *symbolsArchive) StartUpdater(frequency time.Duration) {

	sa.Lock()
	sa.updateStarted = time.Now()
	sa.Unlock()

	go func() {
		updateTicker := time.NewTicker(frequency)
		defer updateTicker.Stop()

		for {
			select {
			case <-sa.stopUpdate:
				log.Printf("Scheduled price update stoppping.")
				return
			case <-updateTicker.C:
				if err := sa.UpdatePrices(); err != nil {
					// log error
					log.Printf("ERROR: updating prices - %s", err)
				}
			}
		}
	}()

}

// Stops automatic price updater
func (sa *symbolsArchive) StopUpdater() {
	sa.Lock()
	sa.stopUpdate <- true
	sa.Unlock()
}

func (sa *symbolsArchive) StartPersistence(dir string) error {
	// check dir exists
	_, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Can't persist to dir: %s - %s", dir, err)
	}

	sa.Lock()
	sa.persistDir = dir
	sa.persistToDisk = true
	sa.Unlock()

	return nil
}

func (sa *symbolsArchive) StopPersistence() {
	sa.Lock()
	sa.persistToDisk = false
	sa.Unlock()
}

func (sa *symbolsArchive) persistPrices(prices []Price) error {

	pricesJSON, err := json.Marshal(&prices)
	if err != nil {
		return err
	}

	// no prices to persist
	if len(prices) == 0 {
		return nil
	}

	// use time of first price in filename

	priceTime := prices[0].At.Format(time.RFC3339)

	priceFilename := priceTime + ".json"

	path := filepath.Join(sa.persistDir, priceFilename)
	return ioutil.WriteFile(path, pricesJSON, os.ModePerm)
}

func (sa *symbolsArchive) LoadPrices(dir string) error {
	// check dir exists
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("Can't load from dir: %s - %s", dir, err)
	}

	// read all files
	for _, file := range files {
		// only load from .json files
		if strings.HasSuffix(file.Name(), ".json") {
			if err := sa.loadPricesFrom(filepath.Join(dir, file.Name())); err != nil {
				return err
			}
		}
	}

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
			fmt.Errorf("Failed to load price from file: %s - %s", filePath, err)
		}
	}

	return nil
}
