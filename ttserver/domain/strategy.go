package domain

import (
	"fmt"
	"sync"
	"time"
)

type Strategy interface {
	ID() string
	Description() string
	ConditionMet(at time.Time) (bool, error)
	CoinPercent() float64
	Symbol() SymbolType
	As() SymbolType
	Start()
	Stop()
	IsRunning() bool
	TriggerCount() int
	IncCount()
}

/*
Strategies are evaluated after every price update to see if their condition(s) have triggered.  Example strategies that "could" be developed are:-

Buy/Sell Strategies
--------------
- (PriceBelowStrategy) price drops below value (v)
- (PriceAboveStrategy) price rises above value (v) (Unlikely use for buying?)
- price drops by (p) percent within (d) duration
- price rises by (p) percent within (d) duration
- price stops rising and starts dropping by (p) percent within (d) duration

*/

type baseStrategy struct {
	sync.RWMutex
	id           string
	symbol       SymbolType
	as           SymbolType
	coinPercent  float64
	isRunning    bool
	triggerCount int
}

// NewBaseStrategy - creates a Base Strategy
func NewBaseStrategy(id string, symbol, as SymbolType, coinPercent float64) (Strategy, error) {

	return newBaseStrategy(id, symbol, as, coinPercent)
}

func newBaseStrategy(id string, symbol, as SymbolType, coinPercent float64) (*baseStrategy, error) {

	if id == "" {
		return nil, fmt.Errorf("ID must be provided")
	}

	if symbol == "" {
		return nil, fmt.Errorf("Symbol must be provided")
	}

	if as == "" {
		return nil, fmt.Errorf("As coin must be provided")
	}

	if coinPercent <= 0 {
		return nil, fmt.Errorf("Coin percentage must be greater than 0")
	}

	if coinPercent > 100 {
		return nil, fmt.Errorf("Coin percentage cannot be greather than 100")
	}

	return &baseStrategy{
		id:          id,
		symbol:      symbol,
		as:          as,
		coinPercent: coinPercent,
	}, nil
}

func (b *baseStrategy) ID() string {
	return b.id
}

func (b *baseStrategy) Symbol() SymbolType {
	return b.symbol
}

func (b *baseStrategy) As() SymbolType {
	return b.as
}

func (b *baseStrategy) CoinPercent() float64 {
	return b.coinPercent
}

// Description - returns a description of the strategy
func (b *baseStrategy) Description() string {
	return "Base strategy for building other strategies upon"
}

// ConditionMet - unsurprisingly does nothing...
func (d *baseStrategy) ConditionMet(at time.Time) (bool, error) {
	return false, nil
}

// IsRunning - is the strategy correctly running?
func (b *baseStrategy) IsRunning() bool {
	b.RLock()
	defer b.RUnlock()
	return b.isRunning
}

func (b *baseStrategy) TriggerCount() int {
	return b.triggerCount
}

func (b *baseStrategy) IncCount() {
	b.triggerCount++
}

// Start - start the strategy running
func (b *baseStrategy) Start() {
	b.Lock()
	b.isRunning = true
	b.Unlock()
}

// Stop - stops the strategy running
func (b *baseStrategy) Stop() {
	b.Lock()
	b.isRunning = false
	b.Unlock()
}

type doNothingStrategy struct {
	baseStrategy
}

// NewDoNothingStrategy - creates a DoNothing Strategy
func NewDoNothingStrategy(id string, symbol, as SymbolType, coinPercentage float64) (Strategy, error) {

	if bs, err := newBaseStrategy(id, symbol, as, coinPercentage); err != nil {
		return nil, err
	} else {
		return &doNothingStrategy{
			baseStrategy: *bs,
		}, nil
	}
}

// Description - returns a description of the strategy
func (d *doNothingStrategy) Description() string {
	return "You say it best... when you say nothing at all..."
}

// ConditionMet - unsurprisingly does nothing...
func (d *doNothingStrategy) ConditionMet(at time.Time) (bool, error) {
	return false, nil
}

type priceAboveStrategy struct {
	baseStrategy
	abovePrice float64
}

// NewPriceAboveStrategy - creates a PriceAbove Strategy
func NewPriceAboveStrategy(id string, symbol, as SymbolType, abovePrice, coinPercentage float64) (Strategy, error) {

	if bs, err := newBaseStrategy(id, symbol, as, coinPercentage); err != nil {
		return nil, err
	} else {
		if abovePrice <= 0 {
			return nil, fmt.Errorf("above price must be greater than 0")
		}
		return &priceAboveStrategy{
			baseStrategy: *bs,
			abovePrice:   abovePrice,
		}, nil
	}
}

// Description - returns a description of the strategy
func (p *priceAboveStrategy) Description() string {
	return fmt.Sprintf("Price Above Strategy\nTriggered when %s price above %f %s - %3.2f%% of coins committed\n", p.symbol, p.abovePrice, p.as, p.coinPercent)
}

// ConditionMet - triggers when price is above
func (p *priceAboveStrategy) ConditionMet(at time.Time) (bool, error) {

	p.RLock()
	defer p.RUnlock()

	if !p.isRunning {
		return false, nil
	}

	// get price at
	price, err := DefaultArchive.GetPriceAs(p.symbol, p.as, at)
	if err != nil {
		return false, fmt.Errorf("Failed to evaluate strategy %s - %s", p.id, err)
	}

	if price.Price > p.abovePrice {
		p.IncCount()
		return true, nil
	}

	return false, nil
}

type priceBelowStrategy struct {
	baseStrategy
	belowPrice float64
}

// NewPriceBelowStrategy - creates a PriceBelow Strategy
func NewPriceBelowStrategy(id string, symbol, as SymbolType, belowPrice, coinPercentage float64) (Strategy, error) {

	if bs, err := newBaseStrategy(id, symbol, as, coinPercentage); err != nil {
		return nil, err
	} else {
		if belowPrice <= 0 {
			return nil, fmt.Errorf("below price must be greater than 0")
		}
		return &priceBelowStrategy{
			baseStrategy: *bs,
			belowPrice:   belowPrice,
		}, nil
	}
}

// Description - returns a description of the strategy
func (p *priceBelowStrategy) Description() string {
	return fmt.Sprintf("Price Below Strategy\nTriggered when %s price below %f %s - %3.2f%% of coins committed\n", p.symbol, p.belowPrice, p.as, p.coinPercent)
}

// ConditionMet - triggers when price is above
func (p *priceBelowStrategy) ConditionMet(at time.Time) (bool, error) {

	p.RLock()
	defer p.RUnlock()

	if !p.isRunning {
		return false, nil
	}

	// get price at
	price, err := DefaultArchive.GetPriceAs(p.symbol, p.as, at)
	if err != nil {
		return false, fmt.Errorf("Failed to evaluate strategy %s - %s", p.id, err)
	}

	if price.Price < p.belowPrice {
		p.IncCount()
		return true, nil
	}

	return false, nil
}
