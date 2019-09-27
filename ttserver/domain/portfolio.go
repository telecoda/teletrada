package domain

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/telecoda/teletrada/exchanges"
	"github.com/telecoda/teletrada/proto"
	"github.com/telecoda/teletrada/ttserver/servertime"
)

type portfolio struct {
	sync.RWMutex
	name     string
	isLive   bool
	balances map[SymbolType]*BalanceAs
}

const DEFAULT_SYMBOL = SymbolType("BTC")

// GetPortfolio returns current portfolio
func (s *server) GetPortfolio(ctx context.Context, req *proto.GetPortfolioRequest) (*proto.GetPortfolioResponse, error) {

	resp := &proto.GetPortfolioResponse{}

	if err := s.updatePortfolios(); err != nil {
		return nil, fmt.Errorf("failed to update portfolios - %s", err)
	}

	balances := s.livePortfolio.balances

	resp.Balances = make([]*proto.Balance, len(balances))

	var err error
	i := 0
	for _, balance := range balances {
		resp.Balances[i], err = balance.toProto()
		if err != nil {
			return nil, err
		}
		i++
	}

	// All balances returned in BTC
	// Calc conversion rate BTC -> as currency

	// convert balances to a different symbol if necessary
	for _, protoBalance := range resp.Balances {
		if protoBalance.As != req.As {
			// convert to different if available

			// TODO - stuff goes here..
		}
	}

	return resp, nil
}

// initPortfolios - fetches latest balances from exchange
func (s *server) initPortfolios() error {

	DefaultLogger.log("Initialising portfolios")
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
			balance.SellStrategy = ss
		}
	}

	if err := s.livePortfolio.reprice(); err != nil {
		return err
	}

	s.simulations = make(map[string]*simulation, 0)
	DefaultLogger.log("Initialised portfolios")
	return nil
}

// updatePortfolios - fetches latest balances and reprices
func (s *server) updatePortfolios() error {
	if err := s.livePortfolio.refreshCoinBalances(); err != nil {
		return err
	}

	if err := s.livePortfolio.reprice(); err != nil {
		return err
	}

	for _, simulation := range s.simulations {
		if simulation.useRealtimeData {
			if err := simulation.reprice(); err != nil {
				return err
			}
		}
	}

	return nil
}

// updateMetrics - sends metrics about portfolios to Influx
func (s *server) saveMetrics() error {

	DefaultLogger.log("Save portfolio metrics")
	// live metrics
	if err := DefaultMetrics.SavePortfolioMetrics(s.livePortfolio); err != nil {
		return err
	}

	// simulated portfolio metrics
	for _, simulation := range s.simulations {
		if err := DefaultMetrics.SavePortfolioMetrics(simulation.portfolio); err != nil {
			return err
		}
	}

	return nil
}

// refreshCoinBalances - fetch latest coin balances from exchange
func (p *portfolio) refreshCoinBalances() error {
	if p == nil {
		return fmt.Errorf("No portfolio to refresh\n")
	}
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
			balance.Total = coinBalance.Free + coinBalance.Locked
			balance.At = servertime.Now()
		} else {
			// new balance
			newBalance := &BalanceAs{
				CoinBalance: coinBalance,
				As:          DEFAULT_SYMBOL,
				At:          servertime.Now(),
				Total:       coinBalance.Free + coinBalance.Locked,
			}
			p.balances[symbol] = newBalance
		}
	}

	return nil
}

// reprice - will reprice all balances based upon latest prices
func (p *portfolio) reprice() error {
	// convert exchange balances to trada balances
	for _, balance := range p.balances {

		if err := balance.reprice(); err != nil {
			return fmt.Errorf("failed reprice balance: %#v - %s", balance, err)
		}
	}
	return nil
}

// repriceAt - will reprice all balances based upon prices at a specific time
func (p *portfolio) repriceAt(at time.Time) error {
	// convert exchange balances to trada balances
	for _, balance := range p.balances {

		if err := balance.repriceAt(at); err != nil {
			return fmt.Errorf("failed reprice balance: %#v - %s", balance, err)
		}

	}
	return nil
}

// repriceAt will reprice balances based upon prices at a specific time
func (b *BalanceAs) repriceAt(at time.Time) error {
	// find latest price for trading pair
	priceAs, err := DefaultArchive.GetPriceAs(SymbolType(b.Symbol), b.As, at)
	if err != nil {
		return fmt.Errorf("failed to get latest price for: %s as %s - %s", b.Symbol, b.As, err)
	}
	return b.repriceUsing(priceAs)
}

// reprice will reprice balances based upon latest prices
func (b *BalanceAs) reprice() error {
	// find latest price for trading pair
	priceAs, err := DefaultArchive.GetLatestPriceAs(SymbolType(b.Symbol), b.As)
	if err != nil {
		return fmt.Errorf("failed to get latest price for: %s as %s - %s", b.Symbol, b.As, err)
	}

	return b.repriceUsing(priceAs)
}

func (b *BalanceAs) repriceUsing(priceAs Price) error {

	if b.Symbol != string(priceAs.Base) {
		return fmt.Errorf("Cannot reprice symbol: %s with price: %s", b.Symbol, priceAs.Base)
	}
	// reprice balance
	b.Price = priceAs.Price
	b.Value = priceAs.Price * b.Total
	b.At = priceAs.At
	b.As = priceAs.As
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
func (p *portfolio) clone() (*portfolio, error) {

	c := &portfolio{
		name:     p.name + "[cloned]",
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

// diff - compares the difference between two portfolios
func (p *portfolio) diff(before *portfolio) (*portfolio, error) {

	d := &portfolio{
		name:     p.name + "[diff]",
		isLive:   false, // diffs are never live
		balances: make(map[SymbolType]*BalanceAs, 0),
	}

	for symbol, balanceNow := range p.balances {

		// fetch before balance
		if balanceBefore, ok := before.balances[symbol]; ok {
			// calc diff
			d.balances[symbol] = &BalanceAs{
				CoinBalance: exchanges.CoinBalance{
					Symbol:   balanceNow.Symbol,
					Exchange: balanceNow.Exchange,
					Free:     balanceNow.Free - balanceBefore.Free,
					Locked:   balanceNow.Locked - balanceBefore.Locked,
				},
				Total:        balanceNow.Total - balanceBefore.Total,
				At:           balanceNow.At,
				As:           balanceNow.As,
				Price:        balanceNow.Price - balanceBefore.Price,
				Value:        balanceNow.Value - balanceBefore.Value,
				Price24H:     balanceNow.Price24H - balanceBefore.Price24H,
				Value24H:     balanceNow.Value24H - balanceBefore.Value24H,
				Change24H:    balanceNow.Change24H - balanceBefore.Change24H,
				ChangePct24H: balanceNow.ChangePct24H - balanceBefore.ChangePct24H,
			}

		} else {
			return nil, fmt.Errorf("Before portfolio did not contain a balance for %s", symbol)
		}
	}

	return d, nil
}

func (p *portfolio) print() {
	fmt.Printf("Portfolio: %s\n", p.name)

	for symbol, balance := range p.balances {
		fmt.Printf("Symbol: %s\n", symbol)
		fmt.Printf("Balance: %#v\n", balance)
	}
}
