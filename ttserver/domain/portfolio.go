package domain

import (
	"fmt"
	"sync"
)

type portfolio struct {
	sync.RWMutex
	balances []*Balance
}

// initPortfolios - fetches latest balances from exchange
func (s *server) initPortfolios() error {
	s.log("Initialising portfolios")
	s.livePortfolio = &portfolio{}
	if err := s.livePortfolio.refreshBalances(); err != nil {
		return err
	}

	s.simPorts = make([]*portfolio, 0)
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (s *server) updatePortfolios() error {
	s.log("Updating portfolios")
	if err := s.livePortfolio.refreshBalances(); err != nil {
		return err
	}
	s.log("Repricing live porfolio")
	if err := s.livePortfolio.reprice(); err != nil {
		return err
	}

	s.log("Repricing simulated porfolios")
	for _, simPort := range s.simPorts {
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

	p.balances = make([]*Balance, 0)

	// convert exchange balances to trada balances
	for _, exchBalance := range exchBalances {

		b := &Balance{
			ExchangeBalance: exchBalance,
			Total:           exchBalance.Free + exchBalance.Locked,
		}
		// lookup symbol
		if symbol, err := DefaultArchive.GetSymbol(SymbolType(exchBalance.Symbol)); err != nil {
			return err
		} else {
			b.Symbol = symbol
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
			return err
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
			fmt.Printf("Exch: %s Sym: %s Free: %f Locked: %f Total: %f USD price: %f USD value %f\n", b.Exchange, b.Symbol.GetType(), b.Free, b.Locked, b.Total, b.LatestUSDPrice, b.LatestUSDValue)
		}
	}
}
