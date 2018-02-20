package domain

import (
	"fmt"
	"sync"
	"time"
)

type portfolio struct {
	sync.RWMutex
	name     string
	isLive   bool
	balances map[SymbolType]*BalanceAs
}

const DEFAULT_SYMBOL = SymbolType("BTC")

// initPortfolios - fetches latest balances from exchange
func (s *server) initPortfolios() error {

	s.log("Initialising portfolios")
	s.livePortfolio = &portfolio{
		name:     "LIVE",
		isLive:   true,
		balances: make(map[SymbolType]*BalanceAs, 0),
	}
	if err := s.livePortfolio.refreshCoinBalances(); err != nil {
		return err
	}

	// set default strategies
	for _, balance := range s.livePortfolio.balances {
		if bs, err := NewBaseStrategy("base-buy", SymbolType(balance.Symbol), balance.As, 100); err != nil {
			return err
		} else {
			balance.BuyStrategy = bs
		}
		if ss, err := NewBaseStrategy("base-sell", SymbolType(balance.Symbol), balance.As, 100); err != nil {
			return err
		} else {
			balance.BuyStrategy = ss
		}
	}

	if err := s.livePortfolio.repriceBalances(); err != nil {
		return err
	}

	s.simPorts = make(map[string]*portfolio, 0)
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (s *server) updatePortfolios() error {
	if err := s.livePortfolio.refreshCoinBalances(); err != nil {
		return err
	}

	if err := s.livePortfolio.repriceBalances(); err != nil {
		return err
	}

	for i, _ := range s.simPorts {
		if err := s.simPorts[i].repriceBalances(); err != nil {
			return err
		}
	}

	return nil
}

// updateMetrics - sends metrics about portfolios to Influx
func (s *server) saveMetrics() error {

	s.RLock()
	defer s.RUnlock()

	// live metrics
	if err := DefaultMetrics.SavePortfolioMetrics(s.livePortfolio); err != nil {
		return err
	}

	// simulated portfolio metrics
	for _, portfolio := range s.simPorts {
		if err := DefaultMetrics.SavePortfolioMetrics(portfolio); err != nil {
			return err
		}
	}

	return nil
}

// refreshCoinBalances - fetch latest coin balances from exchange
func (p *portfolio) refreshCoinBalances() error {
	fmt.Printf("Refreshing balances\n")
	p.Lock()
	defer p.Unlock()

	if !p.isLive {
		return fmt.Errorf("Simulated portfolio: %s cannot have balances refreshed from exchange", p.name)
	}

	coinBalances, err := DefaultClient.GetCoinBalances()
	if err != nil {
		return fmt.Errorf("failed to get balances from exchange: %s", err)
	}

	for _, coinBalance := range coinBalances {
		symbol := SymbolType(coinBalance.Symbol)
		if balance, ok := p.balances[symbol]; ok {
			balance.CoinBalance = coinBalance
		} else {
			// new balance
			newBalance := &BalanceAs{
				CoinBalance: coinBalance,
				As:          DEFAULT_SYMBOL,
				At:          time.Now(),
			}
			p.balances[symbol] = newBalance
		}
	}

	return nil
}

func (p *portfolio) repriceBalances() error {
	// convert exchange balances to trada balances
	for _, balance := range p.balances {

		if err := balance.reprice(); err != nil {
			return fmt.Errorf("failed reprice balance: %#v - %s", balance, err)
		}

	}

	return nil
}

func (b *BalanceAs) reprice() error {
	// find latest price for trading pair
	priceAs, err := DefaultArchive.GetLatestPriceAs(SymbolType(b.Symbol), b.As)
	if err != nil {
		return fmt.Errorf("failed to get latest price for: %s as %s - %s", b.Symbol, b.As, err)
	}

	// reprice balance
	b.Price = priceAs.Price
	b.Value = priceAs.Price * b.Total

	// get 24h price

	daySummary, err := DefaultArchive.GetDaySummaryAs(SymbolType(b.Symbol), b.As)
	if err != nil {
		// no daily price info, but lets carry on
		//return fmt.Errorf("failed to get day summary for: %s as %s - %s", b.Symbol, b.As, err)
	} else {
		b.Price24H = daySummary.ClosePrice
		b.Value24H = b.Price24H * b.Total
		b.Change24H = daySummary.ChangePrice
		b.ChangePct24H = daySummary.ChangePercent
	}

	return nil
}

// clone - creates a clone of portfolio for simulations
func (p *portfolio) clone(newName string) (*portfolio, error) {

	c := &portfolio{
		name:     newName,
		isLive:   false, // clones are never live
		balances: make(map[SymbolType]*BalanceAs, 0),
	}

	for symbol, balance := range p.balances {
		// clone balance
		cb := *balance
		cb.BuyStrategy = nil
		cb.SellStrategy = nil
		c.balances[symbol] = &cb
	}

	return c, nil
}
