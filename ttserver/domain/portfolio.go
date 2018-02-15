package domain

import (
	"fmt"
	"sync"
	"time"
)

type portfolio struct {
	sync.RWMutex
	balances []*BalanceAs
}

const DEFAULT_SYMBOL = SymbolType("BTC")

// initPortfolios - fetches latest balances from exchange
func (s *server) initPortfolios() error {

	s.log("Initialising portfolios")
	s.livePortfolio = &portfolio{}
	if err := s.livePortfolio.refreshBalances(DEFAULT_SYMBOL); err != nil {
		return err
	}

	s.simPorts = make([]*portfolio, 0)
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (s *server) updatePortfolios() error {

	s.log("Updating portfolios")
	if err := s.livePortfolio.refreshBalances(DEFAULT_SYMBOL); err != nil {
		return err
	}
	return nil
}

// refreshBalances - fetch latest balances from exchange
func (p *portfolio) refreshBalances(as SymbolType) error {
	fmt.Printf("Refreshing balances\n")
	p.Lock()
	defer p.Unlock()

	exchBalances, err := DefaultClient.GetBalances()
	if err != nil {
		return fmt.Errorf("failed to get balances from exchange: %s", err)
	}

	p.balances = make([]*BalanceAs, 0)

	// convert exchange balances to trada balances
	for _, exchBalance := range exchBalances {

		b := &BalanceAs{
			Balance: exchBalance,
			Total:   exchBalance.Free + exchBalance.Locked,
			At:      time.Now().UTC(),
			As:      as,
		}

		if err := b.reprice(); err != nil {
			return fmt.Errorf("failed reprice balance: %#v - %s", b, err)
		}

		p.balances = append(p.balances, b)
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
	price24H, err := DefaultClient.GetPriceChange24(b.Symbol, string(b.As))
	if err != nil {
		return fmt.Errorf("failed to get 24h price for: %s as %s - %s", b.Symbol, b.As, err)
	}

	b.Price24H = price24H.Price.Price
	b.Value24H = b.Price24H * b.Total
	b.Change24H = price24H.ChangeAmount
	b.ChangePct24H = price24H.ChangePercent

	return nil
}
