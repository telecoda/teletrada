package domain

import (
	"fmt"
	"log"
	"sync"
	"time"

	client "github.com/influxdata/influxdb/client/v2"
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

	if err := s.livePortfolio.repriceBalances(); err != nil {
		return err
	}

	s.simPorts = make([]*portfolio, 0)
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
	if err := s.livePortfolio.saveMetrics(); err != nil {
		return err
	}

	// simulated portfolio metrics
	for _, portfolio := range s.simPorts {
		if err := portfolio.saveMetrics(); err != nil {
			return err
		}

	}

	return nil
}

func (p *portfolio) saveMetrics() error {

	log.Printf("Sending portfolio balance data to influxdb")
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  DefaultMetrics.dbName,
		Precision: "ns",
	})

	if err != nil {
		return fmt.Errorf("failed to create batch points: %s", err)
	}

	portType := "live"
	if !p.isLive {
		portType = "simulated"
	}

	for _, balance := range p.balances {

		// Create a point and add to batch
		tags := map[string]string{"symbol": string(balance.Symbol), "name": p.name,
			"live": portType}
		fields := make(map[string]interface{}, 0)

		toSymbols := []SymbolType{SymbolType(BTC), SymbolType(ETH), SymbolType(USDT)}

		for _, toSym := range toSymbols {
			if symPrice, err := DefaultArchive.GetLatestPriceAs(SymbolType(balance.Symbol), toSym); err != nil {
				log.Printf("No %s price for %s symbol - %s", toSym, balance.Symbol, err)
			} else {
				fields[fmt.Sprintf("price.%s", toSym)] = symPrice.Price
				// calc current value = total * price
				value := symPrice.Price * balance.Total
				fields[fmt.Sprintf("value.%s", toSym)] = value
			}
		}

		if len(fields) > 0 {
			fields["exchange"] = balance.Exchange
			// add coin totals
			fields["total"] = balance.Total
			fields["locked"] = balance.Locked
			fields["free"] = balance.Free

			// only add fields with points
			pt, err := client.NewPoint("portfolio_balance", tags, fields, balance.At)
			if err != nil {
				fmt.Println("Error: ", err.Error())
			}

			bp.AddPoint(pt)
		}
	}
	// Write the batch
	return DefaultMetrics.Write(bp)

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
func (p *portfolio) clone(newName string) *portfolio {

	c := &portfolio{
		name:     newName,
		isLive:   false, // clones are never live
		balances: make(map[SymbolType]*BalanceAs, 0),
	}

	for symbol, balance := range p.balances {
		// clone balance
		cb := *balance
		c.balances[symbol] = &cb
	}

	return c
}
