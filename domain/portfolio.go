package domain

import (
	"fmt"
	"sync"
)

type portfolio struct {
	sync.RWMutex
	balances []*balance
}

// initPortfolios - fetches latest balances from exchange
func (t *trada) initPortfolios() error {
	fmt.Printf("Initialising portfolios\n")
	t.livePortfolio = &portfolio{}
	if err := t.livePortfolio.refreshBalances(); err != nil {
		return err
	}

	t.simPorts = make([]*portfolio, 0)
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (t *trada) updatePortfolios() error {
	fmt.Printf("Updating portfolios\n")
	if err := t.livePortfolio.refreshBalances(); err != nil {
		return err
	}
	if err := t.livePortfolio.reprice(); err != nil {
		return err
	}

	for _, simPort := range t.simPorts {
		if err := simPort.reprice(); err != nil {
			return err
		}

	}

	return nil
}

// refreshBalances - fetch latest balances from exchange
func (p *portfolio) refreshBalances() error {
	fmt.Printf("Refreshing balances\n")
	p.Lock()
	defer p.Unlock()

	exchBalances, err := DefaultClient.GetBalances()
	if err != nil {
		return err
	}

	p.balances = make([]*balance, 0)

	// convert exchange balances to trada balances
	for _, exchBalance := range exchBalances {

		b := &balance{
			ExchangeBalance: exchBalance,
			Total:           exchBalance.Free + exchBalance.Locked,
		}
		// lookup symbol
		if symbol, err := DefaultArchive.GetSymbol(SymbolType(exchBalance.Symbol)); err != nil {
			return err
		} else {
			b.symbol = symbol
		}

		p.balances = append(p.balances, b)
	}

	return nil
}

// reprice - update balance prices
func (p *portfolio) reprice() error {
	fmt.Printf("Repricing portfolio\n")
	p.Lock()
	defer p.Unlock()

	for _, balance := range p.balances {
		if err := balance.reprice(); err != nil {
			return nil
		}
	}
	return nil
}

func (p *portfolio) ListBalances() {
	fmt.Printf("Balances\n")
	fmt.Printf("========\n")

	p.RLock()
	defer p.RUnlock()
	for _, b := range p.balances {
		if b.Free != 0 {
			fmt.Printf("Exch: %s Sym: %s Free: %f Locked: %f Total: %f USD price: %f USD value %f\n", b.Exchange, b.symbol.GetType(), b.Free, b.Locked, b.Total, b.LatestUSDPrice, b.LatestUSDValue)
		}
	}
}
